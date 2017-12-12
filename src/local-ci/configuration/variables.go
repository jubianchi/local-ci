package configuration

import "github.com/docker/docker/api/types/container"

type Variables struct {
	variables map[string]string
}

func (variables Variables) ConfigureContainer(containerConfiguration *container.Config) {
	for name, value := range variables.variables {
		containerConfiguration.Env = append(containerConfiguration.Env, name+"="+value)
	}
}

func (variables Variables) MergeVariables(right ...Variables) *Variables {
	merged := CreateVariables(make(map[string]string))

	for name, value := range variables.variables {
		merged.variables[name] = value
	}

	for _, vars := range right {
		for name, value := range vars.variables {
			merged.variables[name] = value
		}
	}

	return merged
}

func CreateVariables(variables map[string]string) *Variables {
	return &Variables{
		variables: variables,
	}
}

func CreateVariablesFromYaml(yml map[interface{}]interface{}) *Variables {
	variables := make(map[string]string)

	for name, value := range yml {
		variables[name.(string)] = value.(string)
	}

	return CreateVariables(variables)
}

func ConfigureVariables(configuration map[interface{}]interface{}, variables *Variables) {
	variables.variables["CI_COMMIT_REF_NAME"] = ""
	variables.variables["CI_COMMIT_REF_SLUG"] = ""
	variables.variables["CI_COMMIT_SHA"] = ""
	variables.variables["CI_COMMIT_TAG"] = ""
	variables.variables["CI_DEBUG_TRACE"] = ""
	variables.variables["CI_ENVIRONMENT_NAME"] = ""
	variables.variables["CI_ENVIRONMENT_SLUG"] = ""
	variables.variables["CI_ENVIRONMENT_URL"] = ""
	variables.variables["CI_REPOSITORY_URL"] = ""

	variables.variables["CI_PROJECT_DIR"] = ""
	variables.variables["CI_PROJECT_ID"] = ""

	variables.variables["CI_PROJECT_NAMESPACE"] = ""
	variables.variables["CI_PROJECT_PATH"] = ""
	variables.variables["CI_PROJECT_PATH_SLUG"] = ""
	variables.variables["CI_PROJECT_URL"] = ""
	variables.variables["CI_PROJECT_VISIBILITY"] = ""
	variables.variables["CI_REGISTRY"] = ""
	variables.variables["CI_REGISTRY_IMAGE"] = ""
	variables.variables["CI_REGISTRY_PASSWORD"] = ""
	variables.variables["CI_REGISTRY_USER"] = ""
	variables.variables["CI_SERVER"] = ""
	variables.variables["CI_SERVER_NAME"] = ""
	variables.variables["CI_SERVER_REVISION"] = ""
	variables.variables["CI_SERVER_VERSION"] = ""
	variables.variables["CI_SHARED_ENVIRONMENT"] = ""
	variables.variables["ARTIFACT_DOWNLOAD_ATTEMPTS"] = ""
	variables.variables["GET_SOURCES_ATTEMPTS"] = ""
	variables.variables["GITLAB_CI"] = ""
	variables.variables["GITLAB_USER_ID"] = ""
	variables.variables["GITLAB_USER_EMAIL"] = ""
	variables.variables["GITLAB_USER_LOGIN"] = ""
	variables.variables["GITLAB_USER_NAME"] = ""
	variables.variables["RESTORE_CACHE_ATTEMPTS"] = ""

	for name, value := range configuration {
		variables.variables[name.(string)] = value.(string)
	}
}
