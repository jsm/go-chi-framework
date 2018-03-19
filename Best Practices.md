# Best Practices

## Errors

### Errors of Type String

Always cast an `error` of type `string` to a `string` in the `Error()` function when using an `error` as a parameter in a string format, using `%s` as the format token.

```
type NotFoundError string

func (uuidString NotFoundError) Error() string {
    return fmt.Sprintf("Push Notification not found for uuids: %s", string(uuidString))
}
```
