select Policy {
  name,
  exceptions := .<policies[is Exception] {
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
