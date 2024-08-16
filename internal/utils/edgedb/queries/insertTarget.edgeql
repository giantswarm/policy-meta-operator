with
    existing_target := (
        SELECT Target
        FILTER
            .names = <array<str>>$0 AND
            .namespaces = <array<str>>$1 AND
            .kind = <str>$2
    )
insert Target {
    names := <array<str>>$0,
    namespaces := <array<str>>$1,
    kind := <str>$2
}
unless conflict on (
    .names,
    .namespaces,
    .kind
)
else (
    SELECT existing_target {
        id
    }
);
