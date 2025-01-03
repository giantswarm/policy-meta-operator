insert KyvernoClusterPolicy {
    name := <str>$0,
    ruleNames := <array<str>>$1,
    targetKinds := <array<str>>$2,
    gsExempt := <bool>$3,
}
unless conflict on .name
else (
    update KyvernoClusterPolicy
    set {
        ruleNames := <array<str>>$1,
        targetKinds := <array<str>>$2,
        gsExempt := <bool>$3,
    }
);
