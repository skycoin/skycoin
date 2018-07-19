%extend coin__BlockHeader {
	int __eq__(coin__BlockHeader* a){
		if( $self->Version != a->Version || $self->Time != a->Time || 
			$self->BkSeq != a->BkSeq || $self->Fee != a->Fee)
			return 0;
		if( memcmp( &$self->PrevHash, a->PrevHash, sizeof(a->PrevHash) ) != 0 )
			return 0;
		if( memcmp( &$self->BodyHash, a->PrevHash, sizeof(a->BodyHash) ) != 0 )
			return 0;
		if( memcmp( &$self->UxHash, a->PrevHash, sizeof(a->UxHash) ) != 0 )
			return 0;
		return 1;
	}
}
