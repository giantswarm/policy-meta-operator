create type PolicyException {
  create required property name -> str {
    create constraint exclusive
  };
  create property counter -> int64;
  create property last_reconciliation -> datetime;
}
