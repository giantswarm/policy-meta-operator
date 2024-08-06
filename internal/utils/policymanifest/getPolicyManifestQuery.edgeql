select Policy {
  name,
  automatedExceptions := .<policies[is AutomatedException] {
    targets := (
      select distinct .targets {
        names,
        namespaces,
        id
      }
    )
  },
  policyExceptions := .<policies[is PolicyException] {
    targets := (
      select distinct .targets {
        names,
        namespaces,
        id
      }
    )
  }
}
filter .name = <str>$0;
