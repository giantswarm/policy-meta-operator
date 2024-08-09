insert Policy {
    name := <str>$0,
    defaultPolicyState := <str>$1,
}
unless conflict on .name
else (
    update Policy
    set {
        defaultPolicyState := <str>$1,
    }
);
