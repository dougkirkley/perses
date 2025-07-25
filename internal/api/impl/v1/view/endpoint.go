// Copyright 2024 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package view

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/perses/perses/internal/api/authorization"
	apiInterface "github.com/perses/perses/internal/api/interface"
	"github.com/perses/perses/internal/api/interface/v1/dashboard"
	"github.com/perses/perses/internal/api/interface/v1/view"
	"github.com/perses/perses/internal/api/route"
	"github.com/perses/perses/internal/api/utils"
	v1 "github.com/perses/perses/pkg/model/api/v1"
	"github.com/perses/perses/pkg/model/api/v1/role"
)

type endpoint struct {
	dashboardService dashboard.Service
	service          view.Service
	authz            authorization.Authorization
}

func NewEndpoint(service view.Service, authz authorization.Authorization, dashboardService dashboard.Service) route.Endpoint {
	return &endpoint{
		service:          service,
		authz:            authz,
		dashboardService: dashboardService,
	}
}

func (e *endpoint) CollectRoutes(g *route.Group) {
	g.POST(fmt.Sprintf("/%s", utils.PathView), e.view, false)
}

func (e *endpoint) view(ctx echo.Context) error {
	result := v1.View{}
	if err := ctx.Bind(&result); err != nil {
		return apiInterface.HandleBadRequestError(err.Error())
	}

	if e.authz.IsEnabled() {
		if ok := e.authz.HasPermission(ctx, role.ReadAction, result.Project, role.DashboardScope); !ok {
			return apiInterface.HandleUnauthorizedError(fmt.Sprintf("missing '%s' permission in '%s' project for '%s' kind", role.ReadAction, result.Project, role.DashboardScope))
		}
	}

	if _, err := e.dashboardService.Get(apiInterface.Parameters{
		Project: result.Project,
		Name:    result.Dashboard,
	}); err != nil {
		return apiInterface.HandleNotFoundError(err.Error())
	}

	return e.service.View(&result)
}
