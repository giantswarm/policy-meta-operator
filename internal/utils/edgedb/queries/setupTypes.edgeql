create type Policy {
  create required property name -> str {
    create constraint exclusive
  };
  create required property mode -> str;
  create property last_reconciliation -> datetime;
};

create type Target {
  create required property kind -> str;
  create required property names -> array<str>;
  create required property namespaces -> array<str>;
  create constraint exclusive on ((.kind, .names, .namespaces));
};

create abstract type Exception {
  create multi link targets -> Target;
  create multi link policies -> Policy;
};

create type PolicyException extending Exception {
  create required property name -> str {
    create constraint exclusive
  };
};

create type AutomatedException extending Exception {
  create required property name -> str {
    create constraint exclusive
  };
};
