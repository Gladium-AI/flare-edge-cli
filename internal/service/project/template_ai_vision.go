package project

func aiVisionWorkerMainTemplate() string {
	return aiWorkerTemplate("@cf/moonshotai/kimi-k2.5", `func handleRequest(this js.Value, args []js.Value) any {
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
		prompt = "Describe this image."
	}
	imageURL := searchParam(url, "image_url")
	if imageURL == "" {
		return jsonTextResponse(400, "{\"message\":\"image_url query parameter is required\"}")
	}
	model := searchParam(url, "model")
	if model == "" {
		model = defaultAIModel
	}

	content := js.Global().Get("Array").New()

	textPart := js.Global().Get("Object").New()
	textPart.Set("type", "text")
	textPart.Set("text", prompt)
	content.Call("push", textPart)

	imagePart := js.Global().Get("Object").New()
	imagePart.Set("type", "image_url")
	imageURLPayload := js.Global().Get("Object").New()
	imageURLPayload.Set("url", imageURL)
	imagePart.Set("image_url", imageURLPayload)
	content.Call("push", imagePart)

	messages := js.Global().Get("Array").New()
	user := js.Global().Get("Object").New()
	user.Set("role", "user")
	user.Set("content", content)
	messages.Call("push", user)

	input := js.Global().Get("Object").New()
	input.Set("messages", messages)

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
