package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/service/deployment"
)

// RunDeployment : run a deployment
func RunDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()

	name := params["name"]
	_ = query.Get("tag")

	// Find the deployment
	_, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// jobs.StartDeployment(deployment, tag)

	utils.HTTPAccepted(w)
	return
}

// DeleteDeployment : delete a deployment
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	_, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// jobs.DeleteDeployment(d)

	utils.HTTPAccepted(w)
	return
}

// StopDeployment : stop a deployment
func StopDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	_, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// jobs.StopDeployment(d)

	utils.HTTPAccepted(w)
	return
}

// GetDeployment : get a deployment
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, d)
	return
}

// GetDeployments : get all deployments
func GetDeployments(w http.ResponseWriter, r *http.Request) {
	// Find deployments
	deployments, err := deployment.GetDeployments()
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, deployments)
	return
}
