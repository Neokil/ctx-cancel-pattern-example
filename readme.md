# Context-Cancel-Pattern
In this example I will show how the Context can be used to cancel long running operations. This pattern can be very helpful if you find yourself writing long running functions that need to have a timeout or some other cancel functionality.
The Pattern that we will be using in this example looks like this:
```golang
func foo(ctx context.Context) error {
    // Create one or more result/error channels
    rc := make(chan error, 1)
    defer close(rc)

    go func() {
        err := longRunningBar(ctx)
        if err != nil {
            rc <- err
            return
        }
        rc <- nil
    }()

    select {
        case err := <-rc:
            if err != nil {
                return fmt.Errorf("Make sure to wrap your error: %w", err)
            }
            return nil
        case <- ctx.Done():
            return fmt.Errorf("Make sure to wrap your errors: %w", ctx.Err())
    }
}
```

## Example-Scenario
In the example we have the task to download a list of Users from a Web-API and sync that Data into our Database. Since the WebAPI and the Database Operations can take a long time we want to have a Timeout of 1 Minute for the total Operation, 1s for the API-Calls and 10s for the Database-Queries.
Additionally it would be possible to manually cancel the Context prematurely as I will show.

### Example
`go run main.go api-mock.go db-mock.go`

In the example we have some different timeouts for the contexts and we are using our pattern in all functions.
When a timeout occurs the select-statement will catch it and return an error.
The worker-func in that function will still run, but
1. if it does not handle the cancelation of the context itself it will run through but the result will no longer be handled
2. if it calls more functions that implement the context-cancel-pattern the currently running function will return the wrapped "Deadline Exceeded" Error message (but the worker inside of that function will then run further)
This means even tho we are stopping as soon as the the timeout is reached the previously started worker-routines will continue in the background, so in order to make sure we are having the least amount of overhead due to already canceled operations running we have to make sure our functions are always doing small portions of the work.

## Rollback
If an operation is canceled you might want to do some rollback of the work that was previously done. This is also possible by adding this work to the `case <- ctx.Done():`-Block in the select-statement.
This way you can remove uncomplete data or make sure your data-state is not fucked.

## Takeway
The most important parts of this to take away are:
1. Use the Pattern in all functions that need to be cancelable
2. Make short-running functions to be able to cancel fast
3. Implement rollback to not break your data if necessary
