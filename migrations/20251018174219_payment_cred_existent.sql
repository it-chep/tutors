-- +goose Up
-- +goose StatementBegin
update payment_cred pc
set bank     = 'alpha',
    base_url = 'https://payment.alfabank.ru/payment/rest',
    cred     = json_build_object(
            'user', pc.user_pay,
            'password', pc.password_pay
               );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
