create type Policy {
  create required property name -> str {
    create constraint exclusive
  };
  create multi link exceptions -> AutomatedException;
  create property last_reconciliation -> datetime;
};

create type Target {
  create property names -> array<str>;
  create property namespaces -> array<str>;
  create property kinds -> array<str>
};

create type PolicyException {
  create required property name -> str {
    create constraint exclusive
  };
  create multi link targets -> Target;
  create multi link policies -> Policy;
};

insert Target { 
    names := ["chart-operator"],
    namespaces := ["giantswarm"],
    kinds := ["deployment"]
};
insert Target { 
    names := ["my-deployment"],
    namespaces := ["default"],
    kinds := ["Deployment"]
};