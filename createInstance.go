package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/lestrrat-go/strftime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ec2 "github.com/aws/aws-sdk-go/service/ec2"
)

var ami string     // enter the ami id that you want to use
var size string    // enter the instance size that you wish to use
var sshkey string  // enter name of the ssh-key you wish to use
var subnet string  // enter subnet id that you want to use for instance
var tagName string // enter tag Name value
var tagEnv string  // enter tag environment value prod or staging
var tagApp string  // enter tag app name value
var tagCat string  // enter tag category value

var reader = bufio.NewReader(os.Stdin) // reader for user input from Stdin

const inputdelimiter = '\n' // newline delimiter

// when patching flag "-p" is called this section will run with only input needed is ami.
func patching() {
	now := time.Now()
	t, _ := strftime.Format("%Y-%m-%d", now) // YYYY-MM-DD

	fmt.Print("Please enter ami: ")
	// section for user input from stdin. TO DO clean up to simplify
	ami, err := reader.ReadString(inputdelimiter)
	if err != nil {
		fmt.Println(err)
		return
	}
	ami = strings.Replace(ami, "\n", "", -1) //converts input

	svc := ec2.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-east-1 region
		ImageId:      aws.String(ami),
		InstanceType: aws.String("t2.micro"),
		KeyName:      aws.String("adh-devops"),
		SubnetId:     aws.String("subnet-401ff919"),
		SecurityGroupIds: []*string{
			aws.String("sg-c2f9a2a7"),
			aws.String("sg-8e9dd0eb"),
		},
		//SecurityGroups: []*string{},
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),
	})

	if err != nil {
		log.Println("Could not create instance", err)
		return
	}

	log.Println("Created instance", *runResult.Instances[0].InstanceId)
	// gets username from os
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}
	user := currentUser.Username //user var is set to username

	user = strings.Replace(user, "\n", "", -1)
	tagName = user

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(tagName + "-PATCH"),
			},
			{
				Key:   aws.String("date"),
				Value: aws.String(t),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}
	log.Println("Successfully tagged instance")

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{runResult.Instances[0].InstanceId},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("IP address: ", *inst.PrivateIpAddress)
		}
	}

}

func main() {

	options := flag.Bool("p", false, "p")
	flag.Parse()

	if *options == true {
		patching()
	} else {

		fmt.Print("Please enter ami: ")
		// section for user input from stdin. TO DO clean up to simplify
		ami, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}

		ami = strings.Replace(ami, "\n", "", -1)

		fmt.Print("Please enter size: ")

		size, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}
		size = strings.Replace(size, "\n", "", -1)

		fmt.Print("Please enter sshkey: ")

		sshkey, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}

		sshkey = strings.Replace(sshkey, "\n", "", -1)

		fmt.Print("Please enter subnet: ")

		subnet, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}
		subnet = strings.Replace(subnet, "\n", "", -1)

		svc := ec2.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
		// Specify the details of the instance that you want to create.
		runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
			// An Amazon Linux AMI ID for t2.micro instances in the us-east-1 region
			ImageId:      aws.String(ami),
			InstanceType: aws.String(size),
			KeyName:      aws.String(sshkey),
			SubnetId:     aws.String(subnet),
			SecurityGroupIds: []*string{
				aws.String("sg-c2f9a2a7"),
				aws.String("sg-8e9dd0eb"),
			},
			//SecurityGroups: []*string{},
			MinCount: aws.Int64(1),
			MaxCount: aws.Int64(1),
		})

		if err != nil {
			log.Println("Could not create instance", err)
			return
		}

		log.Println("Created instance", *runResult.Instances[0].InstanceId)

		fmt.Print("Please enter Name tag: ")
		// section for user input from stdin. TO DO clean up to simplify
		tagName, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}
		tagName = strings.Replace(tagName, "\n", "", -1)

		fmt.Print("Please enter app tag: ")

		tagApp, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}
		tagApp = strings.Replace(tagApp, "\n", "", -1)

		fmt.Print("Please enter env tag: ")

		tagEnv, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}

		tagEnv = strings.Replace(tagEnv, "\n", "", -1)

		fmt.Print("Please enter category tag: ")

		tagCat, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			return
		}
		tagCat = strings.Replace(tagCat, "\n", "", -1)

		// Add tags to the created instance
		_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{runResult.Instances[0].InstanceId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(tagName),
				},
				{
					Key:   aws.String("env"),
					Value: aws.String(tagEnv),
				},
				{
					Key:   aws.String("category"),
					Value: aws.String(tagCat),
				},
				/*{
				        Key:   aws.String("date"),
				        Value: aws.String(t),
				},*/
			},
		})

		if errtag != nil {
			log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
			return
		}

		log.Println("Successfully tagged instance")

		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("instance-id"),
					Values: []*string{runResult.Instances[0].InstanceId},
				},
			},
		}
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			fmt.Println("there was an error listing instances in", err.Error())
			log.Fatal(err.Error())
		}

		for idx := range resp.Reservations {
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("IP address: ", *inst.PrivateIpAddress)
			}
		}
	}
}
