for i in {1..100}; do
  valkey-cli -p 6400 LPUSH queue:emails "{\"recipient\":\"test$i@mail.com\",\"template_file\":\"user_templates.tmpl\",\"template_data\":{}}"
done
