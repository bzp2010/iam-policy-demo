# IAM Policy demo

## Models

### User

| Field    | Description               |
| -------- | ------------------------- |
| UserName | string                    |
| Roles    | many2many to Role model   |
| Policies | many2many to Policy model |

### Role

| Field    | Description               |
| -------- | ------------------------- |
| Name     | string                    |
| Users    | many2many to User model   |
| Policies | many2many to Policy model |


### Policy

| Field  | Description                     |
| ------ | ------------------------------- |
| Policy | JSON Policy in PostgreSQL jsonb |
| Users  | many2many to User model         |
| Roles  | many2many to Role model         |

## Test

```bash
curl -s -XPOST "127.0.0.1:8080/hasAccess?username=superadmin"  --data '{"UpdateServiceBasic:A":{"resource":"api:gateway:service:a","action":"gateway:UpdateServiceBasic"},"UpdateServiceBasic:B":{"resource":"api:gateway:service:b","action":"gateway:UpdateServiceBasic"},"EDIT_BUTTON_CLICKABLE":{"resource":"api:gateway:service:a","action":"gateway:UpdateServiceBasic"}}' | jq 

curl -s "127.0.0.1:8080/caniuse_service?username=superadmin" | jq

curl -s "127.0.0.1:8080/caniuse_service?username=alice" | jq

curl "127.0.0.1:8080/hasAccess?username=singlelabel" -s -XPOST --data '{"EDIT_SERVICE_BY_SINGLE_LABEL":{"resource":"api:gateway:service:a","action":"gateway:GetService","sim_labels":{"key:c":"value:c"}}}' | jq # True

curl "127.0.0.1:8080/hasAccess?username=singlelabel" -s -XPOST --data '{"EDIT_SERVICE_BY_SINGLE_LABEL":{"resource":"api:gateway:service:a","action":"gateway:GetService","sim_labels":{"key:d":"value:d"}}}' | jq # False

curl "127.0.0.1:8080/hasAccess?username=singlelabel" -s -XPOST --data '{"EDIT_SERVICE_BY_SINGLE_LABEL":{"resource":"api:gateway:service:a","action":"gateway:GetService","sim_labels":{"key:c":"value:c","key:d":"value:d"}}}' | jq # False

curl "127.0.0.1:8080/hasAccess?username=multiplelabels" -s -XPOST --data '{"EDIT_SERVICE_BY_MUTLIPLE_LABEL":{"resource":"api:gateway:service:a","action":"gateway:GetService","sim_labels":{"key:a":"value:a","key:b":"value:b"}}}' | jq # True

curl "127.0.0.1:8080/hasAccess?username=multiplelabels" -s -XPOST --data '{"EDIT_SERVICE_BY_MUTLIPLE_LABEL":{"resource":"api:gateway:service:a","action":"gateway:GetService","sim_labels":{"key:a":"value:a"}}}' | jq # True

curl "127.0.0.1:8080/hasAccess?username=multiplelabels" -s -XPOST --data '{"EDIT_SERVICE_BY_MUTLIPLE_LABEL":{"resource":"api:gateway:service:a","action":"gateway:GetService","sim_labels":{"key:a":"value:a","key:d":"value:d"}}}' | jq # False
```
