BEGIN;

DELETE FROM
    users
WHERE
    username = 'tester01'
    OR username = 'tester02'
    OR username = 'tester03';

COMMIT;
