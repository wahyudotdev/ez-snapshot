package storage

func New(opts ...Opts) Repository {

	opt := storageOpts{}

	for _, fn := range opts {
		fn(&opt)
	}

	return NewAwsS3Impl()
}
