Releases
========

v1.0.0 (2017-08-29)
-------------------

No changes since v0.3.1. This release is committing to making no breaking
changes to the current API in the 1.X series.


v0.3.1 (2017-05-31)
-------------------

-   Support multierr 1.0.


v0.3.0 (2017-05-24)
-------------------

-   Implement `FieldHook`s natively in mapstructure
-   **Breaking**: Changed function signature of `FieldHook` to remove unecessary
    `from` parameter.

    Before:

    ```go
    func(from reflect.Type, to reflect.StructField, data reflect.Value) (reflect.Value, error)
    ```

    After:

    ```go
    func(dest reflect.StructField, srcData reflect.Value) (reflect.Value, error)
    ```


v0.2.0 (2017-05-03)
-------------------

-   Added `DecodeHook` to intercept values before they are decoded.
-   Added `FieldHook` to intercept values before they are decoded into specific
    struct fields.
-   Decode now parses strings if they are found in place of a float, boolean,
    or integer.
-   Embedded structs and maps will now be inlined into their parent structs.


v0.1.0 (2017-03-31)
-------------------

-   Initial release.
