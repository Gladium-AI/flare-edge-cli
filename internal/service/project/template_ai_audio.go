package project

func aiSTTWorkerMainTemplate() string {
	return aiWorkerTemplate("@cf/deepgram/nova-3", `func handleRequest(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return textResponse(500, "missing request or env")
	}

	request := args[0]
	env := args[1]
	ai := env.Get("AI")
	if ai.IsUndefined() || ai.IsNull() {
		return textResponse(500, "Workers AI binding \"AI\" is not configured")
	}
	if request.Get("body").IsUndefined() || request.Get("body").IsNull() {
		return jsonTextResponse(400, "{\"message\":\"send audio bytes in the request body\"}")
	}

	url := js.Global().Get("URL").New(request.Get("url"))
	model := searchParam(url, "model")
	if model == "" {
		model = defaultAIModel
	}

	contentType := request.Get("headers").Call("get", "content-type")
	if contentType.IsUndefined() || contentType.IsNull() || contentType.String() == "" {
		contentType = js.ValueOf("audio/mpeg")
	}

	audio := js.Global().Get("Object").New()
	audio.Set("body", request.Get("body"))
	audio.Set("contentType", contentType)

	input := js.Global().Get("Object").New()
	input.Set("audio", audio)
	input.Set("detect_language", true)

	options := js.Global().Get("Object").New()
	options.Set("returnRawResponse", true)

	return promise(func(resolve, reject js.Value) {
		var onResolve js.Func
		var onReject js.Func

		onResolve = js.FuncOf(func(this js.Value, args []js.Value) any {
			defer onResolve.Release()
			defer onReject.Release()
			if len(args) > 0 {
				resolve.Invoke(args[0])
				return nil
			}
			resolve.Invoke(jsonTextResponse(502, "{\"message\":\"Workers AI returned no response\"}"))
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

		ai.Call("run", model, input, options).Call("then", onResolve).Call("catch", onReject)
	})
}`)
}

func aiTTSWorkerMainTemplate() string {
	return aiWorkerTemplate("@cf/deepgram/aura-2-en", `func handleRequest(this js.Value, args []js.Value) any {
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
	text := searchParam(url, "text")
	if text == "" {
		text = "Hello from a Go-based Cloudflare text to speech worker."
	}
	model := searchParam(url, "model")
	if model == "" {
		model = defaultAIModel
	}
	speaker := searchParam(url, "speaker")
	if speaker == "" {
		speaker = "luna"
	}

	input := js.Global().Get("Object").New()
	input.Set("text", text)
	input.Set("speaker", speaker)

	options := js.Global().Get("Object").New()
	options.Set("returnRawResponse", true)

	return promise(func(resolve, reject js.Value) {
		var onResolve js.Func
		var onReject js.Func

		onResolve = js.FuncOf(func(this js.Value, args []js.Value) any {
			defer onResolve.Release()
			defer onReject.Release()
			if len(args) > 0 {
				resolve.Invoke(args[0])
				return nil
			}
			resolve.Invoke(jsonTextResponse(502, "{\"message\":\"Workers AI returned no response\"}"))
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

		ai.Call("run", model, input, options).Call("then", onResolve).Call("catch", onReject)
	})
}`)
}
