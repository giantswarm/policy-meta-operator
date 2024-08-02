create type Policy {
  create required property name -> str {
    create constraint exclusive
  };
  create required property mode -> str;
  create property last_reconciliation -> datetime;
};

create type Target {
  create property kind -> str;
  create property names -> array<str>;
  create property namespaces -> array<str>;
  create constraint exclusive on ((.kind, .names, .namespaces));
};

create abstract type Exception {
  create required property name -> str {
    create constraint exclusive
  };
  create multi link targets -> Target;
  create multi link policies -> Policy;
};

create type PolicyException extending Exception {
};

create type AutomatedException extending Exception {
};
