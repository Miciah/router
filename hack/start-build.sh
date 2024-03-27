#!/bin/bash

set -euo pipefail

oc -n openshift-ingress start-build router ${V+--follow} --wait

if [[ -n "${DEPLOY+1}" ]]
then
    oc -n openshift-ingress-operator patch deploy/ingress-operator \
       --type=strategic --patch='
{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "ingress-operator",
            "env": [
              {
                "name": "IMAGE",
                "value": "image-registry.openshift-image-registry.svc:5000/openshift-ingress/origin-haproxy-router:latest"
              }
            ]
          }
        ]
      }
    }
  }
}
'
fi
