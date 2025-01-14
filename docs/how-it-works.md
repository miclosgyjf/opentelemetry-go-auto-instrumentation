# How it works

The workflow could roughly be divided into two main phases:

![](workflow.png)

- `Preprocess`: Analyze dependencies and select rules that should be used later.
- `Instrument`: Generate code based on rules and inject new code into source code.

## Preprocess
Analyze the third-party library dependencies of the user's project code and match 
them with existing instrumentation rules to find suitable rules. Additionally, 
configure the extra dependencies required by these rules in advance. When all 
preprocessing tasks are ready, the `go build -toolexec otel-go-auto-instrumentation cmd/app`
command is invoked for compilation.

Instrumentation rules precisely define which code needs to be injected into which 
version of which framework or standard library. Different types of instrumentation 
rules serve different purposes. The currently available types of instrumentation
rules include:

- InstFuncRule: Inject code at the entry and exit points of a method.
- InstStructRule: Modify a struct by adding a new field.
- InstFileRule: Add a new file to participate in the original compilation process.

The `-toolexec` parameter is the core of automatic instrumentation. It is used to 
intercept the regular build process and replace it with user-defined tools, 
allowing developers more flexibility in customizing the build process. 
The **otel-go-auto-instrumentation** invoked here is the automatic instrumentation 
tool, which leads to the second stage: code injection, i.e. Instrument.

## Instrument
Based on the rules, trampoline code is inserted into the target functions, 
compilation parameters are modified, and then the `go build cmd/app` command is 
invoked for compilation. Trampoline code (Trampoline Jump) is essentially a complex
 If-statement. Through it, instrumentation code can be inserted at the entry and 
 exit points of the target function, enabling the collection of monitoring data.

# `net/http` example
First, we classify the following three types of functions: *RawFunc*, *TrampolineFunc*, *HookFunc*. RawFunc is the original function that needs to be injected. TrampolineFunc is the trampoline function. HookFunc is onEnter/onExit functions that need to be inserted at the entry and exit points of the original function as probe code. RawFunc jumps to TrampolineFunc via the inserted trampoline code, then TrampolineFunc constructs the context, prepares the error recovery handling, and finally jumps to HookFunc to execute the probe code.

![](tjump.png)

Next, we use `net/http` as an example to demonstrate how compile-time automatic instrumentation can insert monitoring code into the target function `(*Transport).RoundTrip()`. The framework will generate trampoline code at the entry of this function, which is an if statement (actually one line, written in multiple lines for demonstration) that jumps to TrampolineFunc:

```go
func (t *Transport) RoundTrip(req *Request) (retVal0 *Response, retVal1 error) {
    if callContext37639, _ := OtelOnEnterTrampoline_RoundTrip37639(&t, &req); false {
    } else {
        defer OtelOnExitTrampoline_RoundTrip37639(callContext37639, &retVal0, &retVal1)
    }
    return t.roundTrip(req)
}
```

Here, `OtelOnEnterTrampoline_RoundTrip37639` is the TrampolineFunc. It prepares error handling and the call context, then jumps to `ClientOnEnterImpl`:

```go
func OtelOnEnterTrampoline_RoundTrip37639(t **Transport, req **Request) (*CallContext, bool) {
    defer func() {
        if err := recover(); err != nil {
            println("failed to exec onEnter hook", "clientOnEnter")
            if e, ok := err.(error); ok {
                println(e.Error())
            }
            fetchStack, printStack := OtelGetStackImpl, OtelPrintStackImpl
            if fetchStack != nil && printStack != nil {
                printStack(fetchStack())
            }
        }
    }()
    callContext := &CallContext{
        Params:     nil,
        ReturnVals: nil,
        SkipCall:   false,
    }
    callContext.Params = []interface{}{t, req}
    ClientOnEnterImpl(callContext, *t, *req)
    return callContext, callContext.SkipCall
}
```

The `ClientOnEnterImpl` is the HookFunc, which is our probe code where traces, metrics reporting, etc., are performed. `ClientOnEnterImpl` is a function pointer, pre-configured in the automatically generated *otel_setup_inst.go* during the preprocessing stage, and it actually points to `clientOnEnter`:

```go
// == otel_setup_inst.go
package otel_rules

import http328 "net/http"
...

func init() {
    http328.ClientOnEnterImpl = clientOnEnter
    ...
}
```

The `clientOnEnter` function performs the actual monitoring tasks:

```go
// == otel_rule_http59729.go
func clientOnEnter(call *http.CallContext, t *http.Transport, req *http.Request) {
    ...
    var tracer trace.Tracer
    if span := trace.SpanFromContext(req.Context()); span.SpanContext().IsValid() {
        tracer = span.TracerProvider().Tracer("")
    } else {
        tracer = otel.GetTracerProvider().Tracer("")
    }
    opts := append([]trace.SpanStartOption{}, trace.WithSpanKind(trace.SpanKindClient))
    ctx, span := tracer.Start(req.Context(), req.URL.Path, opts...)
    var attrs []attribute.KeyValue
    attrs = append(attrs, semconv.HTTPMethodKey.String(req.Method))
    attrs = append(attrs, attributes.MakeSpanAttrs(req.URL.Path, req.URL.Host, attributes.Http)...)
    span.SetAttributes(attrs...)
    bag := baggage.FromContext(ctx)
    if mem, err := baggage.NewMemberRaw(constants.BAGGAGE_PARENT_PID, attributes.Pid); err == nil {
        bag, _ = bag.SetMember(mem)
    }
    if mem, err := baggage.NewMemberRaw(constants.BAGGAGE_PARENT_RPC, sdktrace.GetRpc()); err == nil {
        bag, _ = bag.SetMember(mem)
    }
    sdktrace.SetGLocalData(constants.TRACE_ID, span.SpanContext().TraceID().String())
    sdktrace.SetGLocalData(constants.SPAN_ID, span.SpanContext().SpanID().String())
    ctx = baggage.ContextWithBaggage(ctx, bag)
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
    req = req.WithContext(ctx)
    *(call.Params[1].(**http.Request)) = req
    return
}
```

Through the above steps, we not only inserted monitoring code into the `(*Transport).RoundTrip()` function but also ensured the accuracy and propagation of monitoring data and context. During compile-time automatic instrumentation, these operations are all done automatically, saving developers a significant amount of time and reducing the error rate of manual probes.