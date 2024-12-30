# CLI-GPT

This project is for having a support to have gpt on cli. This is a on going project feel free to push features and issues.

## How to setup the project

- Create a config.yaml file and provide the openai api key, Refer: [Docs](https://platform.openai.com/docs/quickstart/account-setup)
    
        api_key="openai-api-key"

Recent Fixes=
### Changes and Fixes:
1. **Fix config file handling in `initConfig`**

2. **Handle empty choices array properly in `getResFromGpt`**

3. **Refactor error handling and unnecessary `break/os.Exit(1)` in `startChat**

4. **Use `io.ReadAll` instead of `ioutil.ReadAll`**

5. **Fix error handling in `json.Unmarshal`:**

6. **Trim whitespace from user input to check for "exit"**

7. **Provide more informative error messages, especially in API requests"
