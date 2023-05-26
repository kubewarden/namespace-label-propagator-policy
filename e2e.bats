#!/usr/bin/env bats

@test "Accept resource with labels already set" {
  run kwctl run --allow-context-aware -r test_data/pod.json \
	--replay-host-capabilities-interactions test_data/session_replay.yml \
	--settings-path test_data/settings.json annotated-policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
  [ $(expr "$output" : '.*"patchType":"JSONPatch".*') -eq 0 ]
}

@test "Mutate resource with no labels" {
  run kwctl run --allow-context-aware -r test_data/pod_with_no_labels.json \
	--replay-host-capabilities-interactions test_data/session_replay.yml \
	--settings-path test_data/settings.json annotated-policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
  [ $(expr "$output" : '.*"patchType":"JSONPatch".*') -ne 0 ]

  # check if patch operation
  patch=$(echo "${output}" | grep '{.*\"patch\".*}'| jq -r ".patch" | base64 --decode)
  echo "patch=$patch"
  echo ${patch} | jq -e "[.[].op == \"add\" and .[].value.cccenter == \"zpto\"] | any"
}

@test "Mutate resource with wrong labels value" {
  run kwctl run --allow-context-aware -r test_data/pod_with_wrong_label_value.json \
	--replay-host-capabilities-interactions test_data/session_replay.yml \
	--settings-path test_data/settings.json annotated-policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
  [ $(expr "$output" : '.*"patchType":"JSONPatch".*') -ne 0 ]

  # check if patch operation
  patch=$(echo "${output}" | grep '{.*\"patch\".*}'| jq -r ".patch" | base64 --decode)
  echo "patch=$patch"
  echo ${patch} | jq -e "[.[].op == \"replace\" and .[].value == \"zpto\"] | any"
}

@test "Accept resource with labels already set when namespace misses the labels" {
  run kwctl run --allow-context-aware -r test_data/pod.json \
	--replay-host-capabilities-interactions test_data/session_replay_no_labels.yml \
	--settings-path test_data/settings.json annotated-policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "Accept resource with no labels when namespace misses the labels" {
  run kwctl run --allow-context-aware -r test_data/pod_with_no_labels.json \
	--replay-host-capabilities-interactions test_data/session_replay_no_labels.yml \
	--settings-path test_data/settings.json annotated-policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "Accept resource with wrong labels value when namespace misses the labels" {
  run kwctl run --allow-context-aware -r test_data/pod_with_wrong_label_value.json \
	--replay-host-capabilities-interactions test_data/session_replay_no_labels.yml \
	--settings-path test_data/settings.json annotated-policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}
