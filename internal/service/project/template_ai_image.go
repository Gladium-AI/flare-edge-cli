package project

func aiImageWorkerMainTemplate() string {
	return aiWorkerTemplate("@cf/black-forest-labs/flux-2-klein-9b", `func handleRequest(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return textResponse(500, "missing request or env")
	}

	request := args[0]
	env := args[1]
	ai := env.Get("AI")
	if ai.IsUndefined() || ai.IsNull() {
		return textResponse(500, "Workers AI binding \"AI\" is not configured")
	}

	url := js.Global().Get("URL").New(request.Get("url"))
	prompt := searchParam(url, "prompt")
	if prompt == "" {
		prompt = "A dramatic cloudscape over the Alps at sunset."
	}
	model := searchParam(url, "model")
	if model == "" {
		model = defaultAIModel
	}

	form := js.Global().Get("FormData").New()
	form.Call("append", "prompt", prompt)
	form.Call("append", "width", "1024")
	form.Call("append", "height", "1024")

	formResponse := js.Global().Get("Response").New(form)
	multipart := js.Global().Get("Object").New()
	multipart.Set("body", formResponse.Get("body"))
	multipart.Set("contentType", formResponse.Get("headers").Call("get", "content-type"))

	input := js.Global().Get("Object").New()
	input.Set("multipart", multipart)

	return promise(func(resolve, reject js.Value) {
		var onResolve js.Func
		var onReject js.Func

		onResolve = js.FuncOf(func(this js.Value, args []js.Value) any {
			defer onResolve.Release()
			defer onReject.Release()
			payload := "null"
			if len(args) > 0 {
				payload = stringify(args[0])
			}
			resolve.Invoke(jsonResponse(payload))
			return nil
		})

		onReject = js.FuncOf(func(this js.Value, args []js.Value) any {
			defer onResolve.Release()
			defer onReject.Release()
			message := "Workers AI request failed"
			status := 502
			if len(args) > 0 {
				message = errorPayload(args[0])
				if code := errorStatus(args[0]); code > 0 {
					status = code
				}
			}
			resolve.Invoke(jsonTextResponse(status, message))
			return nil
		})

		ai.Call("run", model, input).Call("then", onResolve).Call("catch", onReject)
	})
}`)
}
