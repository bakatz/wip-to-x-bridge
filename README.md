# wip-to-x-bridge
**What**: Finds recently completed todos on WIP.co and posts them to X

**Why**: Zapier doesn't support X

[![Buy me a coffee](https://img.buymeacoffee.com/button-api/?text=Buy%20me%20a%20coffee&emoji=☕&slug=ben_makes_stuff&button_colour=FFDD00&font_colour=000000&font_family=Lato&outline_colour=000000&coffee_colour=ffffff)](https://www.buymeacoffee.com/ben_makes_stuff)

# Requirements (only necessary if you want to build from source, otherwise just skip to the deployment instructions)
- Go 1.x (Navigate to https://go.dev to install the binaries for your OS)
- OSX or Linux (or WSL2 if running on windows)

# Deployment instructions
It's recommended to deploy this to AWS as it's completely free for this use case.
1. Sign up for an AWS account
2. Create a Lambda function, make sure you select `Custom runtime on Amazon Linux 2` as the runtime
3. Change the handler to `bootstrap`
4. Make sure the input is set to Event Bridge with the following settings (change the cron expression depending on how often you want this to run - if you aren't familiar with cron expressions, check out https://crontab.guru):

![system-design](https://github.com/bakatz/rust-server-map-deleter/assets/1575240/3ddaff01-e89e-4094-8a2b-0371dd8f7396)

5. Add the following Environment Variables to the Lambda function you just created:
```
WIP_API_KEY="wipapikey"
TWITTER_API_KEY="twitterapikey"
TWITTER_API_KEY_SECRET="twitterapisecret"
TWITTER_ACCESS_TOKEN="token"
TWITTER_ACCESS_TOKEN_SECRET="tokensecret"
```

6. Go to the latest releases page: https://github.com/bakatz/wip-to-x-bridge/releases and download the lambda-handler.zip file. Alternatively, on your local machine, run ./build.sh which will then output a lambda-handler.zip file.
7. Back in AWS lambda, upload the zip file from the above step under the "Code" menu
8. To test and make sure everything is working, use the Test menu in the AWS Lambda Console to send a test event to the lambda function. It should report back "success." You can also just wait until the scheduled time that you configured as a cron expression and the function will automatically execute.
