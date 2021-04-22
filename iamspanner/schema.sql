CREATE TABLE iam_policy_bindings (
  resource STRING(MAX) NOT NULL,
  binding_index INT64 NOT NULL,
  role STRING(63) NOT NULL,
  member_index INT64 NOT NULL,
  member STRING(MAX) NOT NULL,
) PRIMARY KEY(resource ASC, binding_index ASC, role ASC, member_index ASC, member ASC);

CREATE UNIQUE INDEX iam_policy_bindings_by_member_and_resource ON iam_policy_bindings(member, resource, role);

CREATE UNIQUE INDEX iam_policy_bindings_by_member_and_role ON iam_policy_bindings(member, role, resource);
