// Copyright 2023 The Perses Authors
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

package migrate

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/perses/perses/internal/api/interface"
	"github.com/perses/perses/internal/api/plugin/migrate"
	"github.com/perses/perses/internal/api/route"
	"github.com/perses/perses/pkg/model/api"
)

// Endpoint is the struct that defines all endpoint delivered by the path /migrate
type endpoint struct {
	migrationService migrate.Migration
}

// New create an instance of the object Endpoint.
// You should have at most one instance of this object as it is only used by the struct api in the method api.registerRoute
func New(migrationService migrate.Migration) route.Endpoint {
	return &endpoint{
		migrationService: migrationService,
	}
}

// CollectRoutes is the method to use to register the routes prefixed by /api
// If the version is not v1, then look at the same method but in the package with the version as the name.
func (e *endpoint) CollectRoutes(g *route.Group) {
	g.POST("/migrate", e.Migrate, true)
}

// Migrate is the endpoint that provides the Perses dashboard corresponding to the provided grafana dashboard.
func (e *endpoint) Migrate(ctx echo.Context) error {
	body := &api.Migrate{}
	if err := ctx.Bind(body); err != nil {
		return apiinterface.HandleBadRequestError(err.Error())
	}
	rawGrafanaDashboard := []byte(migrate.ReplaceInputValue(body.Input, string(body.GrafanaDashboard)))
	grafanaDashboard := &migrate.SimplifiedDashboard{}
	if err := json.Unmarshal(rawGrafanaDashboard, grafanaDashboard); err != nil {
		return apiinterface.HandleBadRequestError(err.Error())
	}
	persesDashboard, err := e.migrationService.Migrate(grafanaDashboard)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, persesDashboard)
}
