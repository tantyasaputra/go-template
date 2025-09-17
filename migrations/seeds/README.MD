This directory holds all the seed data scripts needed to have the system up & running.
Scripts are normal goose migrations but should not change the schema.
Scripts should still adhere to reproducibility requirement.

See https://pressly.github.io/goose/blog/2021/no-version-migrations/ for details on how this
differs from normal migrations

To write a seed, the insert must not be duplicated and the must include a down migration that is specific.
Note that the sequence must be maintained between seeds and normal migrations.
This means that a seed with sequence 20 must be followed by sequence 21 regardless of seed or migration.


Example Seed:

```
-- +goose Up
INSERT INTO fdi_ikn_categories
(id, product, provider, provider_detail, payment_channel, fund_pool, created_at, updated_at, deleted_at)
SELECT uuid_generate_v4(), 'TopUp', 'BSI', '7333440444', 'Unique Amount', NULL, 'now()', 'now()', NULL
WHERE
NOT EXISTS (
SELECT * FROM fdi_ikn_categories WHERE provider = 'BSI' and product = 'TopUp' and provider_detail = '7333440444'
);


-- +goose Down
DELETE FROM fdi_ikn_categories WHERE provider = 'BSI' and product = 'TopUp' and provider_detail = '7333440444';


```

Note the insert if not exists pattern