select Policy {
  name,
  mode,
  automatedExceptions := .<policies[is AutomatedException] {
    targets := (
      select distinct .targets {
        kind,
        names,
        namespaces,
        id
      }
    )
  },
  policyExceptions := .<policies[is PolicyException] {
    targets := (
      select distinct .targets {
        kind,
        names,
        namespaces,
        id
      }
    )
  }
}
filter .name = <str>$0;
