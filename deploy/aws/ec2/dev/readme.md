# deployment on EC2 instance.

This deployment is just to create a dev backend in the sky.

## Prerequisistes
A VPC with an internet gateway and public subnet ready to go.

## Setup
- create an EIP for consistent reachability (:money: YOU WILL BE CHARGED AS LONG AS THIS EXISTS)
- had to create an external volume and mount it inside the instance, for pg to use.
- Launch an EC2 instance. A t2.small suffices; could prolly get away with a t2.micro - unsure.
- create a keypair for the instance
- set up a dedicated SG and SG rules:
  - (don't forget will need to access pg, auth and graphql api)
  - ec2 connect ips for setting up SGs (can get them via:)
    ```sh
    curl -s https://ip-ranges.amazonaws.com/ip-ranges.json| jq -r '.prefixes[] | select(.region=="us-east-1") | select(.service=="EC2_INSTANCE_CONNECT") | .ip_prefix'
    ```
    :point_up: whatever ips / cidrs returned need to be allowed for ssh
- Can set env keys directly in the instance and manually run compose; obvs do not do this for anything other than your sandbox. If you will point your local build to the cloud backend, make sure to update your `.env` file.
- can create a launch tpl for simplicity [TODO: add launch tpl or better yet, tf file]