import React from 'react';
import PropTypes from 'prop-types';
import { withNamespaces, Trans } from 'react-i18next';

import Card from '../ui/Card';

const Status = props => (
    <div className="status">
        <Card bodyType="card-body card-body--status">
            <div className="h4 font-weight-light mb-4">
                <Trans>dns_start</Trans>
            </div>
            <button className="btn btn-success" onClick={props.reloadPage}>
                <Trans>try_again</Trans>
            </button>
        </Card>
    </div>
);

Status.propTypes = {
    reloadPage: PropTypes.func.isRequired,
};

export default withNamespaces()(Status);
