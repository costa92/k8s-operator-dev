package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// 过滤所有的PVC,防止其触发 Reconcile 方法
var predicatePVC = predicate.NewPredicateFuncs(func(object client.Object) bool {
	return false
})

var predicatePod = predicate.NewPredicateFuncs(func(object client.Object) bool {
	return false
})
