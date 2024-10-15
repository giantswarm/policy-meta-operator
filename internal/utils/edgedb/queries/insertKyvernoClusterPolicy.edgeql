insert KyvernoClusterPolicy {
    name := <str>$0,
    ruleNames := <str>$1,
    targetKinds := <str>$2,
}
unless conflict on .name
else (
    update KyvernoClusterPolicy
    set {
        ruleNames := <str>$1,
        targetKinds := <str>$2,
    }
);
