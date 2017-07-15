package pipemgr

var JobTemplateDefault = `
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}
  namespace: {{ .Meta.Namespace }}
  labels:
    pipe: {{ .Meta.PipeID }}
    build: {{ .Meta.BuildID }}
    task: {{ .Meta.TaskID }}
spec:
  completions: 1
  parallelism: 1
  template:
    metadata:
      name: {{ .TaskGroup.Title | str2title }}
      labels:
        type: {{ .Meta.Type }}
        pipe: {{ .Meta.PipeID }}
        build: {{ .Meta.BuildID }}
        task: {{ .Meta.TaskID }}
    spec:
{{ range $Index, $Task := .TaskGroup.Tasks }}
      initContainers:
      - name: {{ $Index }}-{{ $Task.Title | str2title }}
        image: {{ $Task.Plugin }}
        imagePullPolicy: {{ $Task.PullPolicy }}
        args:
  {{ range $Task.Args }}
        - {{.}}
  {{ end }}
        command:
  {{ range $Task.Command }}
        - {{.}}
  {{ end }}
        env:
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          value: {{ $Val }}
  {{ end }} 
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          valueFrom:
            secretKeyRef:
              name: {{ $Val }}
              key: data  
  {{ end }}          
{{ end }}
{{ range $Index, $Task := .TaskGroup.Concurrent }}
      containers:
      - name: {{ $Index }}-{{ $Task.Title | str2title }}
        image: {{ $Task.Plugin }}
        imagePullPolicy: {{ $Task.PullPolicy }}
        args:
  {{ range $Task.Args }}
        - {{.}}
  {{ end }}
        command:
  {{ range $Task.Command }}
        - {{.}}
  {{ end }}
        env:
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          value: {{ $Val }}
  {{ end }}
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          valueFrom:
            secretKeyRef:
              name: {{ $Val }}
              key: data  
  {{ end }}
{{ end }}
      restartPolicy: Never

`

var ServiceTemplateDefault = `
apiVersion: batch/v1
kind: Deployment
metadata:
  name: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
  namespace: {{.Meta.Namespace}}
  labels:
    pipe: {{ .Meta.PipeID }}
    build: {{ .Meta.BuildID }}
    task: {{ .Meta.TaskID }}
spec:
  replicas: 1
  template:
    metadata:
      name: {{ .Task.Title | str2title }}
      id: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
      labels:
        pipe: {{ .Meta.PipeID }}
        build: {{ .Meta.BuildID }}
        task: {{ .Meta.TaskID }}
        type: {{ .Meta.Type }}
    spec:
      containers:
      - name: {{ $Index }}-{{ Task.Title | str2title }}
        image: {{ $Task.Plugin }}
        imagePullPolicy: {{ .Task.PullPolicy }}
        ports:
  {{ range .Task.Ports }}
        - containerPort: {{.}}
  {{ end }}
        args:
  {{ range .Task.Args }}
        - {{.}}
  {{ end }}
        command:
  {{ range .Task.Command }}
        - {{.}}
  {{ end }}
        env:
  {{ range $Key, $Val := .Task.Environment }}
        - name: {{ $Key }}
          value: {{ $Val }}
  {{ end }}
  {{ range $Key, $Val := .Task.Environment }}
        - name: {{ $Key }}
          valueFrom:
            secretKeyRef:
              name: {{ $Val }}
              key: data  
  {{ end }}
{{ end }}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
  namespace: {{.Meta.Namespace}}
spec:
  ports:
{{ range $Index, $Val := .Task.Ports }}
    - name: {{ .Index }}
      port: {{ .Val }}
{{ end }}
  selector:
    id: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
`