/*
package config provides support for storing application-wide configuration parameters
in the App Engine Datastore.

To use config, install its middleware in your Kami routes. Then you can retrieve
configuration using config.Get.

You can think of config as an application-wide key-value store. You can store and retrieve
any kind of struct in it that Datastore can serialize. Beware, though, that names will collide
across structs. Consider the following code:

	type AccountConfig struct {
		DefaultRoles []string
	}

	type ActorConfig struct {
		DefaultRoles []string
	}

	actorCfg := &ActorConfig{
		DefaultRoles: []string{"Edward I", "Macbeth"},
	}

	acctCfg := &AccountConfig{
		DefaultRoles: []string{"user", "viewer"},
	}

	config.Save(ctx, actorCfg)
	config.Save(ctx, acctCfg)

In a later request, if somebody does the following,

	type Actor struct {
		Name string
		Role []string
	}

	var currentActorConfig ActorConfig
	config.Get(ctx, &currentActorConfig)

	shakespeare := Actor{
		Name: "William Shakespeare",
		Roles: currentActorConfig.DefaultRoles,
	}


They might be very suprised to find that Bill is set to play "user" and "viewer" rather than "Edward I" and "Macbeth."

*/
package config
