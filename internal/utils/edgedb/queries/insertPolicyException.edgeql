with
    policy_names := <array<str>>$0,
    target_ids := <array<uuid>>$1,
    new_policies := (
      SELECT Policy
      FILTER .name IN array_unpack(policy_names)
    ),
    new_targets := (
      SELECT Target
      FILTER .id IN array_unpack(target_ids)
    ),
insert PolicyException {
    name := <str>$2,
    policies := new_policies,
    targets := new_targets
}
unless conflict on .name
else (
    update PolicyException
    set {
        policies := new_policies,
        targets := new_targets
    }
);
