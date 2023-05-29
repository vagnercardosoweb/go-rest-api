if [ -z "$2" ]; then
  echo "Invalid arguments, use create_template PROFILE_NAME FILE_NAME"
  exit 1
fi

aws --profile "$1" ses update-template --cli-input-json file://"$2".json
