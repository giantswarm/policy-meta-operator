create type Policy {
  create required property name -> str {
    create constraint exclusive
  };
  create property defaultPolicyState -> str;
};

create type PolicyConfig {
  create required property name -> str {
    create constraint exclusive
  };
  create link policyName -> Policy;
  create property policyState -> str;
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

create type KyvernoClusterPolicy {
  create required property name -> str {
    create constraint exclusive
  };
  create required property resourceKinds -> array<str>;
  create required property ruleNames -> array<str>;
}
