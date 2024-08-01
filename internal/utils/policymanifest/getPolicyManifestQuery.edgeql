select Policy {
  name,
  exceptions := .<policies[is Exception] {
    name,
    targets := (
      select distinct .targets {
        names,
        id
      }
    )
  }
}
filter .name = <str>$0;