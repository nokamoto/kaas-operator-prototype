package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"github.com/nokamoto/kaas-operator-prototype/internal/domain"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type objectClient[A client.Object] struct {
	client client.Client
	typ    string
}

func (c *objectClient[A]) create(ctx context.Context, obj A) error {
	if err := c.client.Create(ctx, obj); err != nil {
		return fmt.Errorf("failed to create %s: %w", c.typ, err)
	}
	return nil
}

func (c *objectClient[A]) get(ctx context.Context, name, namespace string, obj A) error {
	err := c.client.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, obj)
	if client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("failed to get %s `%s` in namespace `%s`: %w", c.typ, name, namespace, err)
	}
	if err != nil {
		return errors.Join(domain.ErrResourceNotFound, fmt.Errorf("%s %s not found in namespace %s", c.typ, name, namespace))
	}
	return nil
}
