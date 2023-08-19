package service

import (
	"context"
	"fmt"
	v1 "github.com/costa92/cloud-ide-operator/api/v1"
	"github.com/costa92/cloud-ide-operator/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

const (
	PodNotExist int32 = iota
	PodExit
)

var (
	WorkspaceAlreadyExist = "workspace already exist"
	WorkspaceCreateFailed = "create workspace error"
	WorkspaceNotExist     = "workspace not exist"
	WorkspaceStartFailed  = "start workspace error"
	WorkspaceDeleteFailed = "delete workspace error"
)

var (
	EmptyWorkspaceRunningInfo = &pb.WorkspaceRunningInfo{}
	EmptyResponse             = &pb.Response{}
	EmptyWorkspaceInfo        = &pb.WorkspaceInfo{}
	EmptyWorkspaceStatus      = &pb.WorkspaceStatus{}
)

type WorkSpaceService struct {
	client client.Client
}

func NewWorkSpaceService(c client.Client) *WorkSpaceService {
	return &WorkSpaceService{
		client: c,
	}
}

var _ = pb.CloudIdeServiceServer(&WorkSpaceService{})

// CreateSpace 创建并且启动Workspace,将Operation字段置为"Start",当Workspace被创建时,PVC和Pod也会被创建
func (s *WorkSpaceService) CreateSpace(ctx context.Context, info *pb.WorkspaceInfo) (*pb.WorkspaceRunningInfo, error) {
	// 先查询
	var wp v1.WorkSpace
	exist := s.checkWorkspaceExist(ctx, client.ObjectKey{Name: info.Name, Namespace: info.Namespace}, &wp)
	stus := status.New(codes.AlreadyExists, WorkspaceAlreadyExist)
	if exist {
		return EmptyWorkspaceRunningInfo, stus.Err()
	}
	//不存在就创建
	w := s.constructWorkspace(info)
	if err := s.client.Create(ctx, w); err != nil {
		if errors.IsAlreadyExists(err) {
			return EmptyWorkspaceRunningInfo, stus.Err()
		}
		klog.Errorf("create workspace error:%v", err)
		return EmptyWorkspaceRunningInfo, stus.Err()
	}
	return EmptyWorkspaceRunningInfo, nil
}

func (s *WorkSpaceService) waiteForPodRunning(ctx context.Context, key client.ObjectKey, space v1.WorkSpace) (*pb.WorkspaceRunningInfo, error) {
	// 获取Pod运行的信息。可能会英文资源不足而导致Pod无法运行
	// 最多重试四次，如果还不行，就停止工作空间
	retry, maxRetry := 0, 5
	sleepDuration := []time.Duration{1, 3, 5, 8, 12}
	po := v12.Pod{}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			if retry >= maxRetry {
				break loop
			}
			// 先休眠,等待Pod被创建并且运行起来
			time.Sleep(sleepDuration[retry])

			if err := s.client.Get(context.Background(), key, &po); err != nil {
				if !errors.IsNotFound(err) {
					klog.Errorf("net pod error:%v", err)
				}
			} else {
				if po.Status.Phase == v12.PodRunning {
					return &pb.WorkspaceRunningInfo{
						NodeName: po.Spec.NodeName,
						Ip:       po.Status.PodIP,
						Port:     po.Spec.Containers[0].Ports[0].ContainerPort,
					}, nil
				}
			}
			retry++
		}
	}
	// 5 处理错误情况 停止空间
	if err := s.StopSpace(ctx, &space); err != nil {

	}
	return EmptyWorkspaceRunningInfo, status.New(codes.Internal, WorkspaceStartFailed).Err()
}
func (s *WorkSpaceService) StopSpace(ctx context.Context, space *v1.WorkSpace) (*pb.Response, error) {
	// 使用 Update时 可能由于版本冲突而导致失败，需要重试
	exist := true
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var wp v1.WorkSpace
		exits := s.checkWorkspaceExist(ctx, client.ObjectKey{Name: space.Name, Namespace: space.Namespace}, &wp)
		if !exits {
			return nil
		}

		// 更新Workspace 的Operation 字段
		space.Spec.Operation = v1.WorkSpaceStop
		if err := s.client.Update(ctx, space); err != nil {
			klog.Errorf("update workspace to start error:%v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return EmptyResponse, status.Error(codes.Unknown, err.Error())
	}
	if !exist {
		return EmptyResponse, status.Error(codes.NotFound, WorkspaceNotExist)
	}
	return EmptyResponse, err
}

func (s *WorkSpaceService) checkWorkspaceExist(ctx context.Context, key client.ObjectKey, w *v1.WorkSpace) bool {
	if err := s.client.Get(ctx, key, w); err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		klog.Error("get workspace error:%v", err)
		return false
	}
	return true
}

func (s *WorkSpaceService) constructWorkspace(space *pb.WorkspaceInfo) *v1.WorkSpace {
	hardware := fmt.Sprintf("%sC%s%s", space.ResourceLimit.Cpu,
		strings.Split(space.ResourceLimit.Memory, "i")[0], strings.Split(space.ResourceLimit.Storage, "i")[0])

	return &v1.WorkSpace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cloud-ide.mangohow.com/v1",
			Kind:       "WorkSpace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      space.Name,
			Namespace: space.Namespace,
		},
		Spec: v1.WorkSpaceSpec{
			Cpu:       space.ResourceLimit.Cpu,
			Memory:    space.ResourceLimit.Memory,
			Storage:   space.ResourceLimit.Storage,
			Hardware:  hardware,
			Image:     space.Image,
			Port:      space.Port,
			MountPath: space.VolumeMountPath,
			Operation: v1.WorkSpaceStart,
		},
	}
}
