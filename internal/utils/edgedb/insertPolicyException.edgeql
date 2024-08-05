with
    policy_names := <array<str>>$0,
    new_policies := (
      SELECT Policy
      FILTER .name IN array_unpack(policy_names)
    ),
    targets := (
        insert Target {
            names := <array<str>>$1,
            namespaces := <array<str>>$2,
            kind := <str>$3,
        }
        unless conflict on (
            .names,
            .namespaces,
            .kind
        )
        else (
            select Target
              filter .names = <array<str>>$1
              and .namespaces = <array<str>>$2
              and .kind = <str>$3
            )
    )
insert PolicyException {
    name := <str>$4,
    policies := new_policies,
    targets := targets
}
unless conflict on .name,
else (
    UPDATE PolicyException
    SET {
        policies := new_policies
    }
    FILTER .name = <str>$4
);
