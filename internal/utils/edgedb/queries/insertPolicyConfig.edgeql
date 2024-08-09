insert PolicyConfig {
    name := <str>$0,
    policyState := <str>$1,
    policyName := (
        select Policy
        filter .name = <str>$2
        limit 1
    )
}
unless conflict on .name
else (
    update PolicyConfig
    set {
        policyState := <str>$1,
        policyName := (
            select Policy
            filter .name = <str>$2
            limit 1
        )
    }
);
