# Flows

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/ansible/concord.yml](../concord/examples/ansible/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    # specify the docker image to use
    # add a prefix to use an alternative registry
    dockerImage: "walmartlabs/concord-ansible:latest"
    # rest of the parameters are the usual
    playbook: playbook/hello.yml
    debug: true
    verbose: 3
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/ansible_docker/concord.yml](../concord/examples/ansible_docker/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    dynamicInventoryFile: my_inventory.sh
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/ansible_dynamic_inventory/concord.yml](../concord/examples/ansible_dynamic_inventory/concord.yml)

------

## default

```yaml
# open the form
- form: myForm
# call the playbook
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    extraVars:
        # pass the form field's value into the playbook
        message: ${myForm.myMessage}

```

Defined in [../concord/examples/ansible_form/concord.yml](../concord/examples/ansible_form/concord.yml)

------

## default

```yaml
- form: myForm
  yield: true
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        myHostGroup:
            hosts: ${myForm.ips.split(",")}
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/ansible_form_as_inventory/concord.yml](../concord/examples/ansible_form_as_inventory/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    auth:
        krb5:
            user: "testid"
            password: "PASSWORD"
    inventory:
        myHosts:
            hosts:
                - "testhost"

```

Defined in [../concord/examples/ansible_kerberos/concord.yml](../concord/examples/ansible_kerberos/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
                - "127.0.0.2"
                - "127.0.0.3"
            vars:
                ansible_connection: "local"
    extraVars:
        greetings: "Hi there!"
    limit: "@playbook/hello.limit"

```

Defined in [../concord/examples/ansible_limit/concord.yml](../concord/examples/ansible_limit/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    outVars:
        - "myVar" # created using the `register` statement in the playbook
# `myVar` contains the variable values for all hosts in the play
- log: ${myVar['127.0.0.1']['msg']}

```

Defined in [../concord/examples/ansible_out_vars/concord.yml](../concord/examples/ansible_out_vars/concord.yml)

------

## default

```yaml
# ask the user to fill the form
- form: authForm
  yield: true
- task: ansible
  in:
    # location of the playbook
    playbook: playbook/hello.yml
    # remote server auth
    auth:
        privateKey:
            # remote user's name
            user: "myuser"
            # remote server's key
            secret:
                name: ${authForm.secretName}
                password: ${authForm.password}
    # inventory data, should match the playbook's host groups
    inventory:
        local:
            hosts:
                - "somehost.example.com"
    # pass additional variables to the playbook
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/ansible_remote/concord.yml](../concord/examples/ansible_remote/concord.yml)

------

## default

```yaml
# simply retry immediately after failure
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
                - "127.0.0.2"
                - "127.0.0.3"
            vars:
                ansible_connection: "local"
    extraVars:
        makeItFail: "${makeItFail}"
  retry:
    # specify new task parameters on the retry
    in:
        retry: true # force Ansible to re-use the existing *.retry file
        # this bit is just for example
        extraVars: # override the task's `extraVars` on retry
            makeItFail: false # this time the playbook should succeed
    times: 1
    delay: 3

```

Defined in [../concord/examples/ansible_retry/concord.yml](../concord/examples/ansible_retry/concord.yml)

------

## retryAfterSuspend

```yaml
# Retry the playbook after suspending for some period of time
- try:
    - task: ansible
      in:
        playbook: playbook/hello.yml
        saveRetryFile: true # saves hello.retry as an attachment on playbook error
        limit: "${retryFile}" # default is null, will be a .retry file on retries
        inventory:
            local:
                hosts:
                    - "127.0.0.1"
                    - "127.0.0.2"
                    - "127.0.0.3"
                vars:
                    ansible_connection: "local"
        extraVars:
            makeItFail: "${makeItFail}"
  error:
    - if: "${(attempts + 1) >= maxAttempts}" # give up eventually
      then:
        - throw: "Too many attempts for hosts: ${resource.asString('_attachments/hello.retry')}"
    # suspend until retry time
    - task: sleep
      in:
        suspend: true
        # 3 seconds to prove the point in this example
        # more useful would be 2 hours: ${2 * 60 * 60}
        duration: 3
    # try again, with retry file
    - call: retryAfterSuspend
      in:
        attempts: "${attempts + 1}"
        retryFile: "@_attachments/hello.retry"
        makeItFail: false

```

Defined in [../concord/examples/ansible_retry/concord.yml](../concord/examples/ansible_retry/concord.yml)

------

## retryAfterForm

```yaml
# Prompt to retry playbook after failure
- try:
    - task: ansible
      in:
        playbook: playbook/hello.yml
        saveRetryFile: true # saves hello.retry as an attachment on playbook error
        limit: "${retryFile}" # default is null, will be a .retry file on retries
        inventory:
            local:
                hosts:
                    - "127.0.0.1"
                    - "127.0.0.2"
                    - "127.0.0.3"
                vars:
                    ansible_connection: "local"
        extraVars:
            makeItFail: "${makeItFail}"
  error:
    - if: "${(attempts + 1) >= maxAttempts}" # give up eventually
      then:
        - throw: "Too many attempts for hosts: ${resource.asString('_attachments/hello.retry')}"
    - form: retryForm
      fields:
        - doRetry: {label: "Retry deployment?", type: "boolean"}
        - makeItFail: {label: "Make deployment fail?", type: "boolean"}
    - if: "${retryForm.doRetry}"
      then:
        # try again, with retry file
        - call: retryAfterForm
          in:
            attempts: "${attempts + 1}"
            retryFile: "@_attachments/hello.retry"
            makeItFail: "${retryForm.makeItFail}"
      else:
        - log: "Retry denied"
        - exit

```

Defined in [../concord/examples/ansible_retry/concord.yml](../concord/examples/ansible_retry/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    # location of the playbook
    playbook: "playbook/hello.yml"
    # remote server auth
    auth:
        privateKey:
            # remote user's name
            user: "app"
            # remote server's key
            secret:
                name: "testKey"
    roles:
        - name: "devtools/tekton-ansible"
    inventory:
        myServers:
            hosts:
                - "myRemoteHost"

```

Defined in [../concord/examples/ansible_roles/concord.yml](../concord/examples/ansible_roles/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    extraVars:
        greetings: "Hi there!"
    debug: true
    outVars:
        - "_stats" # register variable for concord ansible stats
- log: "${ _stats }" # will print {failures=[], skipped=[], changed=[], ok=[127.0.0.1], unreachable=[]}

```

Defined in [../concord/examples/ansible_stats/concord.yml](../concord/examples/ansible_stats/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    vaultPassword: myVaultPassword

```

Defined in [../concord/examples/ansible_vault/concord.yml](../concord/examples/ansible_vault/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: playbook/hello.yml
    inventoryFile: inventory.ini
    groupVars:
        - my_hosts:
            secretName: myWindowsKey

```

Defined in [../concord/examples/ansible_windows/concord.yml](../concord/examples/ansible_windows/concord.yml)

------

## default

```yaml
- log: "Starting as ${initiator}"
- form: approvalForm
  runAs:
    ldap:
        group: "CN=Strati-SDE-Concord-sdeconcord,.*"
- if: ${approvalForm.approved}
  then:
    - log: "Approved =)"
  else:
    - log: "Rejected =("

```

Defined in [../concord/examples/approval/concord.yml](../concord/examples/approval/concord.yml)

------

## default

```yaml
- loadTasks: "tasks"
- ${myTask.hey("Concord")}

```

Defined in [../concord/examples/context_injection/concord.yml](../concord/examples/context_injection/concord.yml)

------

## default

```yaml
- form: myForm
  # form calls can override form values or provide additional data
  values:
    lastName: "Appleseed"
    sum: "${1 + 2}"
    address:
        city: Toronto
        province: Ontario
        country: Canada
- log: "Hello, ${myForm.firstName} ${myForm.lastName}"
- log: "We got your file and stored it as ${myForm.aFile}"
- log: "You have following skills"
- task: log
  in:
    msg: "Skill -> ${item}"
  withItems: ${myForm.skills}
- if: ${myForm.tosAgree}
  then:
    - log: "${myForm.firstName} says: I have agreed to the Terms and Conditions"

```

Defined in [../concord/examples/custom_form/concord.yml](../concord/examples/custom_form/concord.yml)

------

## default

```yaml
- form: myForm
  fields:
    - name: {type: "string", value: "${initiator.displayName}"}
- log: "${myForm}"

```

Defined in [../concord/examples/custom_form_basic/concord.yml](../concord/examples/custom_form_basic/concord.yml)

------

## default

```yaml
- task: customTask
  in:
    url: "https://jsonplaceholder.typicode.com/todos/1"
  out: result
- if: "${!result.ok}"
  then:
    - throw: "The request returned ${result.errorCode}"
- log: "Data: ${result.data}"

```

Defined in [../concord/examples/custom_task/test-v2.yml](../concord/examples/custom_task/test-v2.yml)

------

## default

```yaml
- task: customTask
  in:
    url: "https://jsonplaceholder.typicode.com/todos/1"
- if: "${!result.ok}"
  then:
    - throw: "The request returned ${result.errorCode}"
- log: "Data: ${result.data}"

```

Defined in [../concord/examples/custom_task/test.yml](../concord/examples/custom_task/test.yml)

------

## default

```yaml
- log: "${datetime.current()}"
- log: "${datetime.current('dd.MM.yyy')}"
- log: "${datetime.format(datetime.current(), 'dd.MM.yyyy')}"
- log: "${datetime.parse('01.01.2018', 'dd.MM.yyyy')}"

```

Defined in [../concord/examples/datetime/concord.yml](../concord/examples/datetime/concord.yml)

------

## default

```yaml
- docker: "walmartlabs/concord-ansible"
  env:
    ANSIBLE_CONFIG: ansible.cfg
  cmd: ansible-playbook playbook/hello.yml -i inventory.ini -e greetings=Hi

```

Defined in [../concord/examples/docker/concord.yml](../concord/examples/docker/concord.yml)

------

## default

```yaml
- docker: library/alpine
  cmd: echo "Hello, ${name}"

```

Defined in [../concord/examples/docker_simple/concord.yml](../concord/examples/docker_simple/concord.yml)

------

## default

```yaml
# a regular form
- form: myForm1
- log: "${myForm1}"
# an one-off form
- form: myForm2
  fields:
    - firstName: {type: "string", label: "First Name"}
- log: "${myForm2}"
# a form with fields stored in a variable
- form: myForm3
  fields: ${myForm3Fields}
- log: "${myForm3}"
# a form with fields created in a script
- script: groovy
  body: |
    def myFields = [
      ["firstName": ["type": "string", "label": "First Name"]],
      ["lastName": ["type": "string", "label": "Last Name"]]
    ]

    execution.setVariable('myForm4Fields', myFields)
- form: myForm4
  fields: ${myForm4Fields}
- log: "${myForm4}"

```

Defined in [../concord/examples/dynamic_form_fields/concord.yml](../concord/examples/dynamic_form_fields/concord.yml)

------

## default

```yaml
- form: myForm
- log: "I've chosen those colors: ${myForm.colors}"

```

Defined in [../concord/examples/dynamic_form_values/concord.yml](../concord/examples/dynamic_form_values/concord.yml)

------

## default

```yaml
# creating a form using a Groovy script
- script: groovy
  body: |
    // define the form's fields and options
    // the structure is the same as in the form call's syntax
    def myForm = [
      "fields": [
        ["firstName": ["type": "string", "label": "First Name"]],
        ["lastName": ["type": "string", "label": "Last Name"]]
      ],
      "values": [
          "firstName": "John",
          "lastName": "Smith"
      ]
    ]

    // create the form, the process will be suspended after the script is done
    execution.form('myForm', myForm);
- log: "${myForm}"
# creating a form using an expression
- set:
    myForm:
        fields:
            - firstName: {label: "First name", type: "string"}
            - lastName: {label: "Last name", type: "string"}
        values:
            firstName: "John"
            lastName: "Smith"
- ${execution.form('myForm', myForm)}
- log: "${myForm}"

```

Defined in [../concord/examples/dynamic_forms/concord.yml](../concord/examples/dynamic_forms/concord.yml)

------

## default

```yaml
- loadTasks: "tasks"
- ${myTask.hey("world")}

```

Defined in [../concord/examples/dynamic_tasks/concord.yml](../concord/examples/dynamic_tasks/concord.yml)

------

## default

```yaml
- task: loadTasks
  in:
    path: "tasks"
- ${myTask.hey("world")}

```

Defined in [../concord/examples/dynamic_tasks/runtime-v2/concord.yml](../concord/examples/dynamic_tasks/runtime-v2/concord.yml)

------

## default

```yaml
- ${misc.throwBpmnError('kaboom!')}

```

Defined in [../concord/examples/error_handling/concord.yml](../concord/examples/error_handling/concord.yml)

------

## onFailure

```yaml
- log: "Phew! That was close"

```

Defined in [../concord/examples/error_handling/concord.yml](../concord/examples/error_handling/concord.yml)

------

## default

```yaml
- script: example.js
- log: Hello, ${x}

```

Defined in [../concord/examples/external_script/concord.yml](../concord/examples/external_script/concord.yml)

------

## default

```yaml
# "forks" the current process as a child process
- task: concord
  in:
    action: fork
    # if not specified, the parent's entry point will be used
    entryPoint: sayHello
    # wait for completion
    sync: true
    # additional arguments
    arguments:
        otherName: "${initiator.username}"
- log: "Done! ${jobs} is completed"

```

Defined in [../concord/examples/fork/concord.yml](../concord/examples/fork/concord.yml)

------

## sayHello

```yaml
# forked processes can access the latest snapshot of the parent's
# state in addition to the arguments provided by the parent task
- log: "FORK: Hello, ${otherName}. I'm ${myName}"
# simulate a long-running process, sleep for 10s
- ${sleep.ms(10000)}

```

Defined in [../concord/examples/fork/concord.yml](../concord/examples/fork/concord.yml)

------

## default

```yaml
# "forks" the current process as multiple subprocesses
- task: concord
  in:
    action: fork
    tags: forkJoinChild
    # disable the `onCancel` handler, because it's going to handle
    # the parent's cancellation only
    disableOnCancel: true
    forks:
        # spawn multiple jobs with different parameters
        - entryPoint: aJob
          arguments:
            color: "red"
        - entryPoint: aJob
          arguments:
            color: "green"
        - entryPoint: aJob
          arguments:
            color: "blue"
  # out variable "myJobs" will contain a list of process IDs
  out:
    myJobs: ${jobs}
- log: "Done! Status of the jobs: ${concord.waitForCompletion(myJobs)}"

```

Defined in [../concord/examples/fork_join/concord.yml](../concord/examples/fork_join/concord.yml)

------

## aJob

```yaml
- log: "FORK (${color}) starting..."
- ${sleep.ms(15000)}
- log: "...done!"

```

Defined in [../concord/examples/fork_join/concord.yml](../concord/examples/fork_join/concord.yml)

------

## onCancel

```yaml
# find and cancel the tagged subprocesses
- task: concord
  in:
    action: kill
    # because onCancel is a separate subprocess, we need to use
    # 'parentInstanceId'
    instanceId: "${concord.listSubprocesses(parentInstanceId, 'forkJoinChild')}"
    sync: true
- log: "Jobs are cancelled!"

```

Defined in [../concord/examples/fork_join/concord.yml](../concord/examples/fork_join/concord.yml)

------

## default

```yaml
# "yield" makes the process to continue in background after this
# form. It will stop a UI "spinner" in redirects a user back to the
# process' page
- form: myForm
  yield: true
- log: "Hello, ${myForm.name}! I'm starting a long-running task..."
# imitates a long-running task
- ${sleep.ms(30000)}
- log: "Done!"

```

Defined in [../concord/examples/form_and_long_process/concord.yml](../concord/examples/form_and_long_process/concord.yml)

------

## default

```yaml
# this form uses locale.properties in the root directory
- form: myForm
# this form has it's own locale.properties file in forms/myOtherForm directory
# per-form localization files are supported only for custom forms
- form: myOtherForm

```

Defined in [../concord/examples/form_l10n/concord.yml](../concord/examples/form_l10n/concord.yml)

------

## default

```yaml
- form: myForm
- log: "Hello, ${myForm.firstName} ${myForm.lastName}"
- log: "We got your password: ${myForm.password}"
- log: "We know that you are ${myForm.age} years old"
- log: "And your height is: ${myForm.height}"
- log: "'Remember me' checked: ${myForm.rememberMe}"
- log: "File ${myForm.file} content: ${resource.asString(myForm.file)}"
- log: "Email: ${myForm.email}"
- log: "Skills -> ${myForm.skills}"
- task: log
  in:
    msg: "Skill -> ${item}"
  withItems: ${myForm.skills}

```

Defined in [../concord/examples/forms/concord.yml](../concord/examples/forms/concord.yml)

------

## default

```yaml
- log: "Starting as ${initiator}"
- form: myForm
  runAs:
    ldap:
        - group: "CN=Strati-SDE-Concord-sdeconcord,.*"
        - group: "CN=Open Source Developers-opensource_devs,.*"

```

Defined in [../concord/examples/forms_multi_group/concord.yml](../concord/examples/forms_multi_group/concord.yml)

------

## default

```yaml
# call a form using the current (default) value of ${myForm}
- form: myForm
- log: "Value: ${myForm.myField}"
# call a form and override a value
- form: myForm
  values:
    myField: "BBB"
- log: "Value: ${myForm.myField}"
# call a process and override a value in ${myForm}
- call: myFlow
  in:
    myForm:
        myField: "CCC"
- log: "Value: ${myForm.myField}"

```

Defined in [../concord/examples/forms_override/concord.yml](../concord/examples/forms_override/concord.yml)

------

## myFlow

```yaml
# call a form using the current value of ${myForm}
- form: myForm

```

Defined in [../concord/examples/forms_override/concord.yml](../concord/examples/forms_override/concord.yml)

------

## default

```yaml
- call: askUserForDetails

```

Defined in [../concord/examples/forms_wizard/concord.yml](../concord/examples/forms_wizard/concord.yml)

------

## askUserForDetails

```yaml
- form: userData
- if: ${userData.amount > 50}
  then:
    - call: warnUser
  else:
    - call: finishIt

```

Defined in [../concord/examples/forms_wizard/concord.yml](../concord/examples/forms_wizard/concord.yml)

------

## warnUser

```yaml
- form: userWarning
  values:
    amount: ${userData.amount}
- if: ${userWarning.continue == "yes"}
  then:
    - call: finishIt
  else:
    # recursively call the initial form
    - call: askUserForDetails

```

Defined in [../concord/examples/forms_wizard/concord.yml](../concord/examples/forms_wizard/concord.yml)

------

## finishIt

```yaml
- log: "All done!"

```

Defined in [../concord/examples/forms_wizard/concord.yml](../concord/examples/forms_wizard/concord.yml)

------

## default

```yaml
- log: "Running the default flow..."

```

Defined in [../concord/examples/generic_triggers/concord.yml](../concord/examples/generic_triggers/concord.yml)

------

## onEvent

```yaml
- log: "Received ${event}"

```

Defined in [../concord/examples/generic_triggers/concord.yml](../concord/examples/generic_triggers/concord.yml)

------

## onEvent2

```yaml
- log: "${msg}"

```

Defined in [../concord/examples/generic_triggers/concord.yml](../concord/examples/generic_triggers/concord.yml)

------

## default

```yaml
# cloning a repository
- task: git
  in:
    action: clone
    url: git@github.com:myorg/myrepo.git
    privateKey:
        org: myOrg # optional
        secretName: mySecret
        password: myPwd # optional
    workingDir: myRepo
# creating a new branch and pushing it to remote origin
- task: git
  # the repo will be cloned into `myRepo` directory

  in:
    action: createBranch
    url: git@github.com:myorg/myrepo.git
    privateKey:
        org: myOrg # optional
        secretName: mySecret
        password: myPwd # optional
    workingDir: myRepo
    baseBranch: feature-a # optional name of the branch to use as the starting point for the new branch
    newBranch: myNewBranch
    pushBranch: true # set this parameter to 'false' if you do not want to push the new branch to the origin
# merging branches
- task: git
  in:
    action: merge
    url: git@github.com:myorg/myrepo.git
    privateKey:
        org: myOrg # optional
        secretName: mySecret
        password: myPwd # optional
    workingDir: myRepo
    sourceBranch: feature-a
    destinationBranch: myNewBranch
# creating a pull request
- task: github
  in:
    action: createPR
    accessToken: myGitToken
    org: myOrg
    repo: myRepo
    prTitle: "my PullRequest Title"
    prBody: "my PullRequest Body"
    prSourceBranch: mySource # the name of the branch where your changes are implemented.
    prDestinationBranch: master # the name of the branch you want the changes pulled into
# merging a pull request
- task: github
  # the ID of the created PR will be stored as `${prId}`

  in:
    action: mergePR
    accessToken: myGitToken
    org: myOrg
    repo: myRepo
    prId: ${prId}
# create a tag based on a specific commit SHA
- task: github
  in:
    action: createTag
    accessToken: myGitToken
    org: myOrg
    repo: myRepo
    tagVersion: v0.0.1
    tagMessage: "Release 1.0.0"
    tagAuthorName: "myUsedId"
    tagAuthorEmail: "myEmail"
    commitSHA: ${gitHubBranchSHA}

```

Defined in [../concord/examples/git/concord.yml](../concord/examples/git/concord.yml)

------

## default

```yaml
- script: groovy
  body: |
    execution.setVariable("x", 123)
- script: groovy
  body: |
    // variables can be accessed via the context
    println execution.getVariable("x")

    // ...or used directly
    println x

```

Defined in [../concord/examples/groovy/concord.yml](../concord/examples/groovy/concord.yml)

------

## default

```yaml
- script: test.groovy

```

Defined in [../concord/examples/groovy_grape/concord.yml](../concord/examples/groovy_grape/concord.yml)

------

## default

```yaml
- script: groovy
  body: |
    import static groovyx.net.http.HttpBuilder.configure

    def http = configure {
      request.uri = "http://localhost:8001"
    }

    Map result = http.get(Map) {
      request.uri.path = "/api/v1/server/version"
    }

    execution.setVariable("serverVersion", result["version"])
- log: "Server's version: ${serverVersion}"

```

Defined in [../concord/examples/groovy_rest/concord.yml](../concord/examples/groovy_rest/concord.yml)

------

## default

```yaml
- log: "Hello, ${initiator.displayName}"

```

Defined in [../concord/examples/hello_initiator/concord.yml](../concord/examples/hello_initiator/concord.yml)

------

## default

```yaml
- log: "Hello, ${name}!"

```

Defined in [../concord/examples/hello_world/concord.yml](../concord/examples/hello_world/concord.yml)

------

## default

```yaml
- log: "Hello, ${myName}"

```

Defined in [../concord/examples/hello_world2/concord.yml](../concord/examples/hello_world2/concord.yml)

------

## default

```yaml
- log: "${http.asString('http://localhost:8001/api/v1/server/ping')}"
- task: http
  in:
    method: GET
    url: http://localhost:8001/api/v1/server/ping
    response: json
- log: "Response received: ${response}"

```

Defined in [../concord/examples/http/concord.yml](../concord/examples/http/concord.yml)

------

## default

```yaml
# override "name" variable
- call: sayHello
  in:
    name: "${initiator.username}"
# the set value will be kept after the call is ended

# will log "Bye, admin"
- sayBye

```

Defined in [../concord/examples/in_variables/concord.yml](../concord/examples/in_variables/concord.yml)

------

## sayHello

```yaml
- log: "Hello, ${name}"

```

Defined in [../concord/examples/in_variables/concord.yml](../concord/examples/in_variables/concord.yml)

------

## sayBye

```yaml
- log: "Bye, ${name}"

```

Defined in [../concord/examples/in_variables/concord.yml](../concord/examples/in_variables/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    # arguments: [org name], inventory name, ansible host group name, query name, query params, [additional inventory variables]
    inventory: "${inventory.ansible('Default', 'myInventory', 'myHostsGroup', 'endpointsByZypperVersion', {'facter_zypper_version': '1.6.333'})}"
    # will produce a JSON structure like this:
    # {
    #   "myHostsGroup": {
    #     "hosts":["xx.xxx.xx.xxx"]
    #   },
    #   "_meta": {
    #     "hostvars":{
    #       "xx.xxx.xx.xxx":{
    #         "ansible_connection":"local"
    #       }
    #     }
    #   }
    # }

    playbook: playbook/hello.yml
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/inventory/concord.yml](../concord/examples/inventory/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    dockerImage: "walmartlabs/concord-ansible:latest"
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    playbook: playbook/hello.yml
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/inventory_lookup/concord.yml](../concord/examples/inventory_lookup/concord.yml)

------

## default

```yaml
# create new issue
- task: jira
  in:
    action: "createIssue"
    userId: "${initiator.username}"
    password: "myJiraPassword" # should be encrypted and exported using `crypto` task
    projectKey: "MYPROJECTKEY"
    summary: "mySummary"
    description: "myDescription"
    requestorUid: "${initiator.username}"
    issueType: "Bug"
    priority: "P4"
    # complex fields, the actual field names and values depend on the configuration of the JIRA instance
    customFieldsTypeFieldAttr:
        customfield_10212: # environment
            value: "Development"
        customfield_10216: # severity
            value: "4 - Cosmetic"
        customfield_20400: # application/service (array of values)
            - "2125921" # "SDE - Concord"
# the issue ID is stored as `issueId`
- log: "Issue ID: ${issueId}"
# add a comment
- task: jira
  in:
    action: "addComment"
    userId: "${initiator.username}"
    password: "myJiraPassword"
    issueKey: "${issueId}"
    comment: "This is my comment from Concord"
# transition to another status
- task: jira
  in:
    action: "transition"
    userId: "${initiator.username}"
    password: "myJiraPassword"
    issueKey: "${issueId}"
    transitionId: 321
    transitionComment: "Marking as Done"
    customFieldsTypeFieldAttr:
        customfield_10229: # resolution
            value: "Done"
        customfield_20106: # release handling option
            id: "24226" # "This is not going into production (ever)"
# delete an existing issue
- task: jira
  in:
    action: "deleteIssue"
    userId: "${initiator.username}"
    password: "myJiraPassword"
    issueKey: "${issueId}"
- log: "Done!"

```

Defined in [../concord/examples/jira/concord.yml](../concord/examples/jira/concord.yml)

------

## default

```yaml
# using expressions
- log: "${items3.stream().filter(i -> i % 2 == 0).toList()}"
# using Java Streams API in Groovy
- script: groovy
  # calculates the difference between items1 and items2
  body: |
    execution.setVariable("delta",
      items2.stream()
            .filter { a -> !items1.contains(a) }
            .collect())
- log: "We got ${delta}"

```

Defined in [../concord/examples/juel_java_streams/concord.yml](../concord/examples/juel_java_streams/concord.yml)

------

## default

```yaml
- task: ldap
  in:
    action: getUser
    ldapAdServer: "ldap://ldaphost:port"
    bindUserDn: "myBindUser"
    bindPassword: "myBindPwd"
    searchBase: "DC=MyOrg,DC=com"
    user: "myUser"
- log: "${ldapResult}"

```

Defined in [../concord/examples/ldap/concord.yml](../concord/examples/ldap/concord.yml)

------

## default

```yaml
- script: groovy
  body: |
    import org.slf4j.*

    def logger = org.slf4j.LoggerFactory.getLogger("test")
    logger.debug("Hi there")

```

Defined in [../concord/examples/logback_config/concord.yml](../concord/examples/logback_config/concord.yml)

------

## default

```yaml
- ${log.debug("This is a debug log")}
- ${log.info("This is a info log")}
- ${log.warn("This is a warn log")}
- ${log.error("This is a error log")}
- logDebug: "Hello, ${name}. This is a debug log"
- log: "Hello, ${name}. This is a normal log"
- logWarn: "Hello, ${name}. This is a warn log"
- logError: "Hello, ${name}. This is a error log"

```

Defined in [../concord/examples/loglevel/concord.yml](../concord/examples/loglevel/concord.yml)

------

## default

```yaml
- log: "Taking a nap..."
- ${sleep.ms(60000)}
- log: "Done napping!"

```

Defined in [../concord/examples/long_running/concord.yml](../concord/examples/long_running/concord.yml)

------

## default

```yaml
# calling a task for each element
- task: log
  in:
    msg: ${item}
  withItems:
    - "Hello!"
    - "Bye!"
# calling a flow for each element
- call: myFlow
  withItems:
    - "first element"
    - "second element"
# using a variable
- call: myFlow
  withItems: ${myItems}

```

Defined in [../concord/examples/loops/concord.yml](../concord/examples/loops/concord.yml)

------

## myFlow

```yaml
- log: "We got ${item}"

```

Defined in [../concord/examples/loops/concord.yml](../concord/examples/loops/concord.yml)

------

## default

```yaml
- call: myFlow

```

Defined in [../concord/examples/mocking/concord.yml](../concord/examples/mocking/concord.yml)

------

## myFlow
the flow we want to test

```yaml
# normally, it would call the "real" task
- task: github
  in:
    action: mergePR
    accessToken: "..."
    org: myOrg
    repo: myRepo
    prId: 123

```

Defined in [../concord/examples/mocking/concord.yml](../concord/examples/mocking/concord.yml)

------

## test

```yaml
- loadTasks: "mocks"
- try:
    - call: myFlow
  error:
    - log: "Failed as expected with ${lastError.cause}"

```

Defined in [../concord/examples/mocking/concord.yml](../concord/examples/mocking/concord.yml)

------

## default

```yaml
- script: javascript
  body: |
    execution.setVariable("name", "Concord");
# call another flow
- sayHello
# call another flow and add/override its variables
- call: sayHello
  in:
    name: "world"

```

Defined in [../concord/examples/multiple_flows/concord.yml](../concord/examples/multiple_flows/concord.yml)

------

## sayHello

```yaml
- log: "Hello, ${name}"

```

Defined in [../concord/examples/multiple_flows/concord.yml](../concord/examples/multiple_flows/concord.yml)

------

## default

```yaml
- call: findHostsWithArtifacts
- call: findFacts
- call: findDeployedOnHosts

```

Defined in [../concord/examples/noderoster/concord.yml](../concord/examples/noderoster/concord.yml)

------

## findHostsWithArtifacts

```yaml
- task: noderoster
  in:
    action: "hostsWithArtifacts"
    artifactPattern: "storesystems"
- log: "${result}"

```

Defined in [../concord/examples/noderoster/concord.yml](../concord/examples/noderoster/concord.yml)

------

## findFacts

```yaml
- task: noderoster
  in:
    action: "facts"
    hostName: "host.example.com"
- log: "${result}"

```

Defined in [../concord/examples/noderoster/concord.yml](../concord/examples/noderoster/concord.yml)

------

## findDeployedOnHosts

```yaml
- task: noderoster
  in:
    action: "deployedOnHost"
    hostName: "host.example.com"
- log: "${result}"

```

Defined in [../concord/examples/noderoster/concord.yml](../concord/examples/noderoster/concord.yml)

------

## default

```yaml
- set:
    x: 123
    y:
        some:
            nested: ["data", "in", "arrays"]
            boolean: true
            number: 234

```

Defined in [../concord/examples/out/concord.yml](../concord/examples/out/concord.yml)

------

## default

```yaml
- script: groovy
  body: |
    execution.setVariable("myVar", "myValue");

```

Defined in [../concord/examples/out_groovy/concord.yml](../concord/examples/out_groovy/concord.yml)

------

## default

```yaml
# read a JSON file
# prints out '{value:${nested.value}}'
- log: "${resource.asJson('my.json')}"
# read a JSON file and evaluate all expressions
# prints out '{value:Hello!}'
- log: "${resource.asJson('my.json', true)}"
# read a YAML file
# prints out '{value:${nested.value}}'
- log: "${resource.asYaml('my.yml')}"
# read a YAML file and evaluate all expressions
# prints out '{value:Hello!}'
- log: "${resource.asYaml('my.yml', true)}"

```

Defined in [../concord/examples/parsing_yaml_json/concord.yml](../concord/examples/parsing_yaml_json/concord.yml)

------

## default

```yaml
# executes the provided payload archive as a child process
- task: concord
  in:
    action: start
    archive: example.zip
    # wait for completion
    sync: true
- log: "Done! ${jobs[0]} is completed"

```

Defined in [../concord/examples/process_from_a_process/concord.yml](../concord/examples/process_from_a_process/concord.yml)

------

## default

```yaml
# starts a process and retrieves it's variable value
- task: concord
  in:
    action: start
    archive: out.zip
    sync: true
    # list of output variables
    outVars:
        - result
- log: "${jobOut.result}"
# starts a process from an existing project (it must be created beforehand)
- task: concord
  in:
    action: start
    project: test
    repository: default
    sync: true
- log: "Done! ${jobs[0]} is completed"

```

Defined in [../concord/examples/process_from_a_process2/concord.yml](../concord/examples/process_from_a_process2/concord.yml)

------

## default

```yaml
# uses the specified directory as the process payload
- task: concord
  in:
    action: start
    payload: example
    arguments:
        name: "Concord"
    sync: true
- log: "Done! ${jobs[0]} is completed"

```

Defined in [../concord/examples/process_from_a_process3/concord.yml](../concord/examples/process_from_a_process3/concord.yml)

------

## default

```yaml
# use a local file and a process argument to create the message
- log: "${resource.asString('file.txt')}, ${name}!"

```

Defined in [../concord/examples/process_from_a_process3/example/concord.yml](../concord/examples/process_from_a_process3/example/concord.yml)

------

## default

```yaml
- log: "Hello, ${name}!"
- form: myForm
- log: "Action: ${myForm.action}!"

```

Defined in [../concord/examples/process_meta/concord.yml](../concord/examples/process_meta/concord.yml)

------

## default

```yaml
- log: "Hello, ${name}!"

```

Defined in [../concord/examples/profiles/concord/concord.yml](../concord/examples/profiles/concord/concord.yml)

------

## default

```yaml
- log: "Hello, ${name}"

```

Defined in [../concord/examples/profiles/concord.yml](../concord/examples/profiles/concord.yml)

------

## default

```yaml
- log: Hello, ${name}

```

Defined in [../concord/examples/project_file/concord.yml](../concord/examples/project_file/concord.yml)

------

## default

```yaml
- script: example.py
- log: "The result: ${y}"

```

Defined in [../concord/examples/python_script/concord.yml](../concord/examples/python_script/concord.yml)

------

## default

```yaml
- script: ruby
  body: |
    puts "Hello!"
    $execution.setVariable("test", "foo");
- log: "done: ${test}"

```

Defined in [../concord/examples/ruby/concord.yml](../concord/examples/ruby/concord.yml)

------

## default

```yaml
- set:
    url: https://concord.prod.walmart.com
    # in v2, variables are scoped to a flow, unlike global variables in v1
- task: http
  in:
    url: "${url}" # required to be passed. in v1, this was implicit
    method: "GET"
  out: response
- log: "Response Code: ${response.statusCode}"
  # Task inputs are explicit in v2 – all required parameters must be specified in the `in` block
- name: Log Me! # log segments can be named
  task: log
  in:
    msg: "Hello! I'm being logged in a separate (and named!) segment!"
    level: "WARN"
    # separate and customizable log segments allow cleaner readability and log management
- call: flow-v2

```

Defined in [../concord/examples/runtime-v2/a_basic_example/concord/example.concord.yml](../concord/examples/runtime-v2/a_basic_example/concord/example.concord.yml)

------

## flow-v2

```yaml
- try:
    - log: "${invalid expression}"
  error:
    - log: "${lastError}" # will print error log

```

Defined in [../concord/examples/runtime-v2/a_basic_example/concord/example.concord.yml](../concord/examples/runtime-v2/a_basic_example/concord/example.concord.yml)

------

## default

```yaml
- task: ansible
  in:
    playbook: "playbook/hello.yml"
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    outVars:
        - "myVar" # created using the `register` statement in the playbook
  out: ansibleResult
# `myVar` contains the variable values for all hosts in the play
- log: "${ansibleResult.myVar['127.0.0.1']['msg']}"

```

Defined in [../concord/examples/runtime-v2/ansible_out_vars/concord.yml](../concord/examples/runtime-v2/ansible_out_vars/concord.yml)

------

## formFlow
easier to access form links in v2
execute forms in parallel

```yaml
- log: "Going to do some form stuff..."
- parallel:
    - block:
        - form: getFullName
        - log: "${getFullName.firstName}"
    - block:
        - form: fetchMetadata
        - log: "${fetchMetadata.age}"
    - block:
        - form: approvalCRQ
        - log: "${approvalCRQ.password}"

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord/forms.concord.yml](../concord/examples/runtime-v2/demo-flow/concord/forms.concord.yml)

------

## externalFlow

```yaml
# named log segments
- name: ${item}
  task: log
  in:
    msg: "Hello! I'm being logged in a separate segment named ${item}!"
  withItems:
    - "Run Validation"
    - "Trigger Deployment"
    - "Send Notification"
    # Other UI/UX features:
# parallel execution
# The v2 runtime was designed with parallel execution in mind. It adds a new step - parallel.
- parallel:
    - ${sleep.ms(5000)}
    - ${sleep.ms(5000)}
  # show effective concord yaml
  # Show task debug params by recording in/out vars
  # View yaml line number in the Events tab
- log: "Total sleeping duration should be 5 seconds!"
  # executes each expression in its own Java thread
# parallel execution in loop
- task: http
  in:
    url: https://concord.${item}.walmart.com/
    method: "GET"
    debug: true
  out: results # withItems turns "results" into a list of results for each item
  parallelWithItems:
    - "prod"
    - "test"
    - "ci"
- log: ${results.stream().map(o -> o.statusCode).toList()}

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord/test.concord.yml](../concord/examples/runtime-v2/demo-flow/concord/test.concord.yml)

------

## default

```yaml
# The "allVariables()" function returns a Java Map with all the provided process variables
# The "hasVariable()" function accepts a string parameter and returns true if the variable exists
- log: "All variables: ${allVariables()}"
- if: ${hasVariable('projectInfo.orgName')}
  then:
    - log: "Yep, we have 'orgName' variable with value ${projectInfo.orgName}"
  else:
    - log: "Nope, we do not have 'orgName' variable"
- call: externalFlow
- call: scriptFlow
- call: formFlow
- call: errorFlow

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord.yml](../concord/examples/runtime-v2/demo-flow/concord.yml)

------

## aFlow
variable scoping changes

```yaml
- set:
    x: 123
- log: "${x}" # prints out "123"
- call: anotherFlow
- log: "${x}" # prints out "123"
- call: anotherFlow
  out: x # required if output needs to be used in the parent flow
- log: "${x}" # prints out "789"
- call: taskFlow

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord.yml](../concord/examples/runtime-v2/demo-flow/concord.yml)

------

## anotherFlow

```yaml
- log: "${x}" # prints out "123"
- set:
    x: 789

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord.yml](../concord/examples/runtime-v2/demo-flow/concord.yml)

------

## taskFlow
explicit task inputs

```yaml
- set:
    url: https://github.com
- name: "Concord Endpoint Call"
  task: http
  in:
    url: "${url}" # will not use the global variable as tasks now use local variables
    method: "GET" # need to explicitly specify each input parameter
    debug: true
  out: response # need to explicitly specify output var
- if: ${response.ok} # in v2, http response object contains the `ok` attribute instead of the `success` attribute in v1
  then:
    - log: "Concord Endpoint Status: ${response.statusCode}"
  else:
    - log: "Something went wrong!"
- call: default

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord.yml](../concord/examples/runtime-v2/demo-flow/concord.yml)

------

## scriptFlow
get and set script variables

```yaml
- script: scripts/test-script.groovy
  out: newVar
  meta:
    segmentName: "Processing ${myVar} Script Variable..."
  error:
    - log: "${lastError}" # prints original exception, nothing added by concord
- log: "The new value is: ${newVar}" # prints new value set inside the groovy script

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord.yml](../concord/examples/runtime-v2/demo-flow/concord.yml)

------

## errorFlow
improved error syntax

```yaml
- log: "uncomment the lines to see improved error messages!"

```

Defined in [../concord/examples/runtime-v2/demo-flow/concord.yml](../concord/examples/runtime-v2/demo-flow/concord.yml)

------

## default

```yaml
- script: groovy
  body: |
    result.set("myVar", "myValue");
  out: scriptResult
- log: "result: ${scriptResult}" # result: {myVar=myValue}
- script: groovy
  body: |
    result.set("myVar", "myValue");
  out:
    myVar: ${result.myVar}
- log: "myVar: ${myVar}" # myVar: myValue

```

Defined in [../concord/examples/runtime-v2/out_groovy/concord.yml](../concord/examples/runtime-v2/out_groovy/concord.yml)

------

## default

```yaml
- script: js
  body: |
    result.set("myVar", "myValue");
  out: scriptResult
- log: "result: ${scriptResult}" # result: {myVar=myValue}
- script: js
  body: |
    result.set("myVar", "myValue");
  out:
    myVar: ${result.myVar}
- log: "myVar: ${myVar}" # myVar: myValue

```

Defined in [../concord/examples/runtime-v2/out_js/concord.yml](../concord/examples/runtime-v2/out_js/concord.yml)

------

## default

```yaml
- script: python
  body: result.set("myVar", "myValue");
  out: scriptResult
- log: "result: ${scriptResult}" # result: {myVar=myValue}
- script: python
  body: result.set("myVar", "myValue");
  out:
    myVar: ${result.myVar}
- log: "myVar: ${myVar}" # myVar: myValue

```

Defined in [../concord/examples/runtime-v2/out_python/concord.yml](../concord/examples/runtime-v2/out_python/concord.yml)

------

## default

```yaml
- script: ruby
  body: |
    $result.set("myVar", "myValue");
  out: scriptResult
- log: "result: ${scriptResult}" # result: {myVar=myValue}
- script: ruby
  body: |
    $result.set("myVar", "myValue");
  out:
    myVar: ${result.myVar}
- log: "myVar: ${myVar}" # myVar: myValue

```

Defined in [../concord/examples/runtime-v2/out_ruby/concord.yml](../concord/examples/runtime-v2/out_ruby/concord.yml)

------

## default

```yaml
# The v1 runtime provides no satisfactory way to run flow steps in parallel in one single process.
# The v2 runtime was designed with parallel execution in mind. It adds a new step - parallel.
- parallel:
    - ${sleep.ms(5000)}
    - ${sleep.ms(5000)}
- log: "Total sleeping duration should be 5 seconds!"
# sequence of tasks can be run inside the block statement of the parallel step
- parallel:
    - block:
        - name: "HTTP Google Segment"
          task: http
          in:
            url: https://google.com/
            method: "GET"
          out: googleResponse
        - log: "Google: ${googleResponse.statusCode}"
    - block:
        - name: "HTTP Bing Segment"
          task: http
          in:
            url: https://bing.com/
            method: "GET"
          out: bingResponse
        - log: "Bing: ${bingResponse.statusCode}"
  out:
    - googleResponse
    - bingResponse
- log: |-
    Google: ${googleResponse}
    Bing: ${bingResponse}

```

Defined in [../concord/examples/runtime-v2/parallel_execution/concord.yml](../concord/examples/runtime-v2/parallel_execution/concord.yml)

------

## default

```yaml
# the script file has to be served by some external web server
- script: "http://localhost:8000/example.groovy"

```

Defined in [../concord/examples/script_url/concord.yml](../concord/examples/script_url/concord.yml)

------

## default

```yaml
# exporting secrets as files
- set:
    myFileA: ${crypto.exportAsFile('myFileA', pwd)}
    myFileB: ${crypto.exportAsFile('myFileB', pwd)}
# the resulting variables will contain the path of the exported files
- log: "My file A: ${myFileA}"
- log: "My file B: ${myFileB}"

```

Defined in [../concord/examples/secret_files/concord.yml](../concord/examples/secret_files/concord.yml)

------

## default

```yaml
- task: ansible
  in:
    dockerImage: "walmartlabs/concord-ansible:latest"
    playbook: playbook/hello.yml
    inventory:
        local:
            hosts:
                - "127.0.0.1"
            vars:
                ansible_connection: "local"
    extraVars:
        greetings: "Hi there!"

```

Defined in [../concord/examples/secret_lookup/concord.yml](../concord/examples/secret_lookup/concord.yml)

------

## default

```yaml
# working with SSH key pairs
- set:
    myKey: ${crypto.exportKeyAsFile('myKey', pwd)}
- log: "Public key file: ${myKey.public}"
- log: "Private key file: ${myKey.private}"
# working with username/password credentials
- log: "Credentials: ${crypto.exportCredentials('myCreds', pwd)}"
# working with encrypted strings
- log: "Plain secret: ${crypto.exportAsString('myValue', pwd)}"

```

Defined in [../concord/examples/secrets/concord.yml](../concord/examples/secrets/concord.yml)

------

## default

```yaml
- task: slack
  in:
    channelId: "C5W9ELY7Q"
    text: "Process ${txId} is completed"
    username: "my bot"
    iconEmoji: ":chart_with_upwards_trend:"
    attachments:
        - fallback: "Book your flights at https://flights.example.com/book/r123456"
          actions:
            - type: "button"
              text: "Book flights"
              url: "https://flights.example.com/book/r123456"
- log: notified

```

Defined in [../concord/examples/slack/concord.yml](../concord/examples/slack/concord.yml)

------

## default

```yaml
# create a new slack channel
- task: slackChannel
  in:
    action: create
    channelName: myChannelName
    apiToken: mySlackApiToken
# the channel ID is stored as `slackChannelId`
- log: "Channel ID: ${slackChannelId}"
# archive a slack channel
- task: slackChannel
  in:
    action: archive
    channelId: ${slackChannelId}
    apiToken: mySlackApiToken
# create a new slack group
- task: slackChannel
  in:
    action: createGroup
    channelName: myGroupName
    apiToken: mySlackApiToken
# the channel ID is stored as `slackChannelId`
- log: "Channel ID: ${slackChannelId}"
# archive a slack group
- task: slackChannel
  in:
    action: archiveGroup
    channelId: ${slackChannelId}
    apiToken: mySlackApiToken

```

Defined in [../concord/examples/slackChannel/concord.yml](../concord/examples/slackChannel/concord.yml)

------

## default

```yaml
- task: smtp
  in:
    # a custom SMTP server can be specified here.
    # Otherwise, the task will use the global SMTP configuration.

    #smtpParams:
    #  host: "localhost"
    #  port: 25
    mail:
        from: ${initiator.attributes.mail}
        to: ${initiator.attributes.mail}
        subject: "Howdy!"
        template: "mail.mustache"
        attachments:
            - "first.txt"
            - path: "second.txt"
              disposition: "attachment" # optional, valid values: "attachment" or "inline"
              description: "attachment description" # optional
              name: "attachment name" # optional
- log: mail sent to ${initiator.attributes.mail}

```

Defined in [../concord/examples/smtp/concord.yml](../concord/examples/smtp/concord.yml)

------

## default

```yaml
- task: smtp
  in:
    # a custom SMTP server can be specified here.
    # Otherwise, the task will use the global SMTP configuration.

    #smtpParams:
    #  host: "localhost"
    #  port: 25
    mail:
        from: ${initiator.attributes.mail}
        to: ${initiator.attributes.mail}
        subject: "Howdy!"
        template: "mail.mustache.html"
- log: mail sent to ${initiator.attributes.mail}

```

Defined in [../concord/examples/smtp_html/concord.yml](../concord/examples/smtp_html/concord.yml)

------


