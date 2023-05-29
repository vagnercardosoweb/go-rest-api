if [ -z "$2" ]; then
  echo "Invalid arguments, use create_template PROFILE_NAME TEMPLATE_NAME"
  exit 1
fi

aws --profile "$1" ses delete-template --template-name "$2"
