/*
package slug provides support for "slug" keys, which are string IDs for Datastore keys that
need to be usable as URL components.

Let's say you want to key the Profile entity on a URL-friendly form of a person's name. You setttle on something like:

  /profiles/john-smith

Suppose, as can happen, you have two different John Smiths in your system. Keying on name might be a problem
unless you can autoincrement the slug thus:

  /profiles/john-smith-2

Which is precisely the behavior permitted by slug.Next. You can write code for this that looks like:

  nextSlug, err := slug.Next(ctx, "Profile", slugify(newProfile.name))
  datastore.Put(ctx, datastore.NewKey(ctx, "Profile", nextSlug, 0, nil), &newProfile)

*/
package slug
