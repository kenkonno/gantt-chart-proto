resource "aws_vpc" "test_vpc" {
}

resource "aws_ecs_cluster" "test_ecs_cluster" {
  name = ""
}

resource "aws_autoscaling_group" "test_aws_autoscaling_group" {

  max_size = 0
  min_size = 0
}


resource "aws_launch_configuration" "api_launch_configuration" {

  image_id      = ""
  instance_type = ""
}
resource "aws_launch_template" "api_aws_launch_template" {
}