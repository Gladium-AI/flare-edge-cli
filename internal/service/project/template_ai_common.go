package project

import (
	"strconv"
	"strings"
)

func aiWorkerTemplate(defaultModel, handleRequest string) string {
	parts := []string{
		"package main\n\n",
		"import \"syscall/js\"\n\n",
		"const defaultAIModel = " + strconv.Quote(defaultModel) + "\n\n",
		"func main() {\n",
		"\tjs.Global().Set(\"handleRequest\", js.FuncOf(handleRequest))\n",
		"\tjs.Global().Set(\"handleScheduled\", js.FuncOf(func(this js.Value, args []js.Value) any {\n",
		"\t\treturn nil\n",
		"\t}))\n",
		"\tselect {}\n",
		"}\n\n",
		handleRequest,
		"\n\n",
		aiWorkerCommonRuntime,
	}
	return strings.Join(parts, "")
}

const aiWorkerCommonRuntime = `func promise(fn func(resolve js.Value, reject js.Value)) js.Value {
	executor := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) >= 2 {
			fn(args[0], args[1])
		}
		return nil
	})
	defer executor.Release()

	return js.Global().Get("Promise").New(executor)
}

func stringify(value js.Value) string {
	output := js.Global().Get("JSON").Call("stringify", value)
	if output.IsUndefined() || output.IsNull() {
		return "null"
	}
	return output.String()
}

func errorPayload(value js.Value) string {
	payload := js.Global().Get("Object").New()
	payload.Set("name", value.Get("name"))
	payload.Set("message", value.Get("message"))
	payload.Set("code", firstDefined(value, "code", "status"))
	payload.Set("internalCode", firstDefined(value, "internalCode", "internal_code"))

	cause := value.Get("cause")
	if !cause.IsUndefined() && !cause.IsNull() {
		causePayload := js.Global().Get("Object").New()
		causePayload.Set("name", cause.Get("name"))
		causePayload.Set("message", cause.Get("message"))
		causePayload.Set("code", firstDefined(cause, "code", "status"))
		causePayload.Set("internalCode", firstDefined(cause, "internalCode", "internal_code"))
		payload.Set("cause", causePayload)
	}

	text := stringify(payload)
	if text == "{}" || text == "null" {
		return "{\"message\":\"Workers AI request failed\"}"
	}
	return text
}

func errorStatus(value js.Value) int {
	for _, key := range []string{"status", "statusCode", "code"} {
		field := value.Get(key)
		if !field.IsUndefined() && !field.IsNull() && field.Type() == js.TypeNumber {
			return field.Int()
		}
	}
	return 0
}

func firstDefined(value js.Value, keys ...string) js.Value {
	for _, key := range keys {
		field := value.Get(key)
		if !field.IsUndefined() && !field.IsNull() {
			return field
		}
	}
	return js.Null()
}

func searchParam(url js.Value, key string) string {
	value := url.Get("searchParams").Call("get", key)
	if value.IsUndefined() || value.IsNull() {
		return ""
	}
	return value.String()
}

func jsonResponse(body string) js.Value {
	headers := js.Global().Get("Headers").New()
	headers.Call("set", "content-type", "application/json")
	options := js.Global().Get("Object").New()
	options.Set("status", 200)
	options.Set("headers", headers)
	return js.Global().Get("Response").New(body, options)
}

func jsonTextResponse(status int, body string) js.Value {
	headers := js.Global().Get("Headers").New()
	headers.Call("set", "content-type", "application/json")
	options := js.Global().Get("Object").New()
	options.Set("status", status)
	options.Set("headers", headers)
	return js.Global().Get("Response").New(body, options)
}

func textResponse(status int, body string) js.Value {
	headers := js.Global().Get("Headers").New()
	headers.Call("set", "content-type", "text/plain; charset=utf-8")
	options := js.Global().Get("Object").New()
	options.Set("status", status)
	options.Set("headers", headers)
	return js.Global().Get("Response").New(body, options)
}
`
