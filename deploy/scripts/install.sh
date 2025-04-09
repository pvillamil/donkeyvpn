#!/bin/bash
set -e


function wizard() {
    go run cmd/installer/main.go \
        -tfvars-template "$PWD/deploy/terraform/terraform.tfvars.example" \
        -tfvars-output "$PWD/deploy/terraform/terraform.tfvars" \
        -tfbackend-template "$PWD/deploy/terraform/terraform.tfbackend.example" \
        -tfbackend-output "$PWD/deploy/terraform/terraform.tfbackend"
}

function build() {
    echo
    echo
    echo "Building Donkeyvpn binary..."
    go build -o dist/bootstrap cmd/bot/main.go
    echo "DonkeyVPN built successfully"
}

function terraform_apply() {
    echo
    echo
    echo "Starting Terraform plan and apply"
    pushd $PWD/deploy/terraform
    terraform init
    terraform apply
    popd
    echo
    echo
    echo "Installation finished successfully! Send /help to your telegram bot"

    # This file is generated by terraform
    if [[ -f /tmp/donkeyvpn-webhook-register.sh ]]; then
        cat /tmp/donkeyvpn-webhook-register.sh
        NO_FORMAT="\033[0m"
        C_GREY0="\033[38;5;16m"
        C_SEAGREEN1="\033[48;5;84m"
        echo -e "${C_GREY0}${C_SEAGREEN1}Run the previous commands ro the next script to register the bot webhook${NO_FORMAT}"
        echo -e "${C_GREY0}${C_SEAGREEN1}bash /tmp/donkeyvpn-webhook-register.sh ${NO_FORMAT}"
    fi
}

function all() {
    wizard
    build
    terraform_apply
}

func=${1:-all}
$func
