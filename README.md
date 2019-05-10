Plentiform
===

An open source form backend for your JAMstack sites. https://plentiform.com

## Developer Guide

### Prerequisites

1. Install [golang](https://golang.org/)

<details>
<summary>How to install golang-1.12 on Ubuntu 18.04</summary>

```shell
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install golang-go
sudo cp -R /usr/lib/go-1.12 /usr/local/go
```
Add the following to ~/.bashrc:
```shell
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```
Then `source ~/.bashrc` to apply changes.

Check that it worked: `go version` should return something like "go version go1.12 linux/amd64"

More info: https://github.com/golang/go/wiki/Ubuntu

</details>

2. Install [docker](https://www.docker.com/)

<details>
<summary>How to install docker on Ubuntu 18.04</summary>

```shell
sudo apt update
sudo apt install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
sudo apt update
sudo apt install docker-ce
```
To be able to run `docker` without `sudo`:
```shell
sudo usermod -aG docker ${USER}
su - ${USER}
```

Check that it worked: `docker -v` should return something like "Docker version 18.09.5, build e8ff056"

More info: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-18-04

</details>

3. Install [docker-compose](https://github.com/docker/compose)

<details>
<summary>How to install docker-compose-1.24.0 on Ubuntu 18.04</summary>

```shell
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

Check that it worked: `docker-compose --version` should return something like "docker-compose version 1.24.0, build 0aa59064"

More info: https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-18-04

</details>

4. Install [migrate](https://github.com/golang-migrate/migrate)

<details>
<summary>How to install golang-migrate on Ubuntu 18.04</summary>

```shell
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ bionic main" | sudo tee -a /etc/apt/sources.list.d/migrate.list
sudo apt-get update
sudo apt-get install -y migrate
go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
```

Check that it worked: `migrate -version` should return something like "4.3.1"

More info: https://github.com/golang-migrate/migrate/tree/master/cli

</details>

### Installation

1. Setup environment variables for things like the database connection, the port your app runs on, and reCAPTCHA
```shell
cp .env.example .env
```

2. Start up a docker container for our Postgres database
```shell
docker-compose up -d
```

3. Start the built in webserver so you can see the app in your browser at http://localhost:3000/
```shell
go run .
```

4. Run the migrations to create tables in the database
```shell
migrate -path migrations -database 'postgres://postgres@localhost:5433?sslmode=disable' up
```

5. Create a default user so you can login (email=admin@example.com, password=admin)
```shell
docker exec -it plentiform_db_1 psql -U postgres -c "insert into users (name, email, password_digest, is_email_confirmed) values ('admin', lower('admin@example.com'), crypt('admin', gen_salt('bf', 8)), true) returning *;"
```

### Creating Forms

You can create a new form right browser by clicking the "forms" link in the main nav (http://localhost:3000/forms).
Then click the blue "Create Form" button (http://localhost:3000/forms/new).
For now just name your form, the "Description" and "ReCaptcha Secret Key" are optional.
Once created, you'll be redirected back to the Overview of all your forms (http://localhost:3000/forms).
You should be able to see your new form in the list.
The form's UUID will be displayed in a format similar to: `98011ea6-f0d5-4c33-b3e5-84b61c3cf8f5`
Copy this, because we'll use it in our example for submitting forms.

### Submitting Forms

To test form submissions, we need to make an HTTP POST request.
One of the easiest ways to do this is to install [Insomnia](https://github.com/getinsomnia/insomnia) (On Ubuntu 18.04: `sudo snap install insomnia`).

Create a new POST request with the following URL (substitute the UUID with the one from your form and make sure your app is running):
```
http://localhost:3000/f/98011ea6-f0d5-4c33-b3e5-84b61c3cf8f5
```
In the Insomnia app where it says "body" click the carrot icon to the right and select "Structured > Multipart Form".
Add whatever form values you'd like to submit. For example:
- "name": "Jimbo"
- "hobby": "forms"

Finally press "send" in the Insomnia app and go back to Plentiform to check the submissions for the form (http://localhost:3000/forms/1).

### reCAPTCHA

Create a reCAPTCHA here (Needed to choose v2): https://www.google.com/recaptcha
Add the site key and secret key to your local .env file

### Accessing the Database Directly

- Open a Postgres prompt inside your database container: `docker-compose exec -u postgres db psql`
- List all databases: `\l`
- Switch to the "postgres" database we created: `\c postgres`
- List the tables in our database: `\dt`
- Create a new user with email=test@example.com and password=test: `insert into users (name, email, password_digest, is_email_confirmed) values ('test', lower('test@example.com'), crypt('test', gen_salt('bf', 8)), true) returning *;`
- Show all users in the "users" table: `select * from users;`
- Get the UUID from your first form in the "forms" table: `select uuid from forms where id=1;`
- Exit Postgres prompt: `\q`

## Credit

This project is a fork of [LetterDrop](https://github.com/jonahgeorge/letterdrop) by [Jonah George](https://twitter.com/jonahgeorge). Big thanks to him sharing his code and [helping us](https://github.com/jonahgeorge/letterdrop/issues/20) get started.
