insert Target {
    names := <array<str>>$0,
    namespaces := <array<str>>$1,
    kind := <str>$2,
}
unless conflict on (
    .names,
    .namespaces,
    .kind
);
